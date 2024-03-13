package command

import (
	"bufio"
	"context"
	"fmt"
	"homework2/pup/cmd/parsing"
	"homework2/pup/internal/model"
	"homework2/pup/internal/service"
	"os"
	"strconv"
	"sync"
)

// help prints usage guide
func Help() {
	fmt.Println(`
	usage: go run ./cmd -command=<help|accept|remove|give|list|return|list-return|pickpoints> [-id=<order id>] [-recipient=<recipient id>] [-expire=<expire date>] [-t=<bool>] [<args>]

	Command desciption:
		help: список доступных команд с кратким описанием
		accept: принять заказ от курьера
		remove: вернуть заказ курьеру
		give: выдать заказ клиенту
		list: получить список заказов
		return: принять возврат от клиента
		list-return: получить список возвратов
		pickpoints: активация интерактивного режима записи и чтения данных о ПВЗ

	Needed flags or arguments for each command:
		help	
		accept 		 -id -recipient -expire
		remove  	 -id
		give		 args: order ids to give (example: go run ./cmd -command=give 1 2 3 4)
		list		 -recipient (optional flag -t: boolean value for printing orders located in our point (not already given); optional args: number of orders to list or zero for all)
		return  	 -id -recipient
		list-return	 args: page number and number of orders per page (default: all pages and 10 orders per page) (example: "-command=list-return 2 5" prints 2nd page of returned orders grouped by 5 orders in each page)
		pickpoints
	
	Flags requirements:
		-id, -recipient: positive number
		-expire: date in 'dd.mm.yyyy' format (02.01.2006 for 2nd Jan 2006)
	`)
	HelpPickPoints()
}

func Accept(serv service.Service, params parsing.Params) {
	if params.ID == nil || params.RecipientID == nil || params.ExpireString == nil {
		fmt.Println("miss required flags")
		return
	}
	err := serv.AcceptFromCourier(model.OrderInput{
		ID:          *params.ID,
		RecipientID: *params.RecipientID,
		ExpireDate:  *params.ExpireString,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("got new order from courier")
}

func Remove(serv service.Service, params parsing.Params) {
	if params.ID == nil || params.RecipientID == nil || params.ExpireString == nil {
		fmt.Println("miss required flags")
		return
	}
	err := serv.Remove(*params.ID)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("removed order %d from our pick-up point\n", *params.ID)
}

func Give(serv service.Service, params parsing.Params) {
	if len(params.Args) == 0 {
		fmt.Println("expected at least one argument as order id")
		return
	}
	ids := make([]int64, len(params.Args))
	for i, s := range params.Args {
		idCur, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			fmt.Println(err)
			return
		}
		ids[i] = idCur
	}
	err := serv.Give(ids)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("all orders have been given to its recipient")
}

func List(serv service.Service, params parsing.Params) {
	if params.RecipientID == nil {
		fmt.Println("miss required flags")
		return
	}
	var (
		n   int
		err error
	)
	if len(params.Args) > 0 {
		n, err = strconv.Atoi(params.Args[0])
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	foundArr, err := serv.List(*params.RecipientID, n, *params.NotGiven)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("found %d orders:\n", len(foundArr))
	for i, order := range foundArr {
		fmt.Printf("%d.\tid: %d\texpires: %s\n", i+1, order.ID, order.ExpireDate.Format("01.02.2006"))
	}
}

func Return(serv service.Service, params parsing.Params) {
	if params.ID == nil || params.RecipientID == nil {
		fmt.Println("miss required flags")
		return
	}
	err := serv.Return(*params.ID, *params.RecipientID)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("order %d is returned successfully\n", *params.ID)
}

func ListReturn(serv service.Service, params parsing.Params) {
	var (
		pageNum, ordersPerPage int
		err                    error
	)
	ordersPerPage = 10
	if len(params.Args) > 0 {
		pageNum, err = strconv.Atoi(params.Args[0])
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	if len(params.Args) > 1 {
		ordersPerPage, err = strconv.Atoi(params.Args[1])
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	arr, err := serv.ListReturn(pageNum, ordersPerPage)
	if err != nil {
		fmt.Println(err)
		return
	}
	startPos := 1
	if pageNum == 0 {
		fmt.Println("all returned not removed orders:")
	} else {
		startPos = ordersPerPage*(pageNum-1) + 1
		fmt.Printf("returned not removed orders from page %d (%d-%d):\n", pageNum, startPos, startPos+len(arr)-1)
	}

	for i, order := range arr {
		fmt.Printf("%d.\tid: %d\trecipient: %d\texpires: %s\n", startPos+i, order.ID, order.RecipientID, order.ExpireDate.Format("01.02.2006"))
	}
}

const chanSize = 10

func PickPoints(serv service.Service) {
	var (
		line, com string
		id        int64
		point     model.PickPoint
		wg        sync.WaitGroup
	)
	ctx, cancel := context.WithCancel(context.Background())
	wg.Add(2)
	writeChan := make(chan model.PickPoint, chanSize)
	readChan := make(chan int64, chanSize)
	go serv.WritePoints(ctx, writeChan, &wg)
	go serv.ReadPoints(ctx, readChan, &wg)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line = scanner.Text()
		fmt.Sscanf(line, "%s", &com)
		switch com {
		case "help":
			HelpPickPoints()
		case "exit":
			cancel()
			wg.Wait()
			return
		case "write":
			_, err := fmt.Sscanf(line, "write %d %s %s %s", &point.ID, &point.Name, &point.Address, &point.Contact)
			if err != nil {
				fmt.Println(err)
				continue
			}
			writeChan <- point
			fmt.Println(point)
		case "read":
			_, err := fmt.Sscanf(line, "read %d", &id)
			if err != nil {
				fmt.Println(err)
				continue
			}
			readChan <- id
			fmt.Println(id)
		default:
			fmt.Println("Unknown command")
		}
	}
}

func HelpPickPoints() {
	fmt.Println(`
	interactive mode for command pickpoints usage guide:

	Command desciption:
		help: список доступных команд с кратким описанием
		write: добавить информацию о ПВЗ
		read: считать информацию о ПВЗ
		exit: завершение работы

	Needed arguments for each command:
		help	
		write 		 id(int)	name(string)	address(string)	   contact(string)
		read	  	 id(int)
		exit
	
	Examples:
		write 10 Chertanovo Chertanovskaya-Street-10 +7(999)888-77-66
		read 10
	`)
}
