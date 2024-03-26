package command

import (
	"fmt"
	"strconv"

	"gitlab.ozon.dev/mer_marat/homework/cmd/console-app/parsing"
	"gitlab.ozon.dev/mer_marat/homework/internal/model"
)

type serviceOrder interface {
	AcceptFromCourier(model.OrderInput) error
	Remove(int64) error
	Give([]int64) error
	List(int64, int, bool) ([]model.Order, error)
	Return(int64, int64) error
	ListReturn(int, int) ([]model.Order, error)
}

func Help() {
	fmt.Println(`
	usage: go run ./cmd -command=<help|accept|remove|give|list|return|list-return|pickpoints> [-id=<order id>] [-recipient=<recipient id>] [-weight=<order weight>] [-price=<order price>] [-cover=<order cover>] [-expire=<expire date>] [-t=<bool>] [<args>]

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
		accept 		 -id -recipient -weight -price -cover -expire
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

func Accept(serv serviceOrder, params parsing.Params) {
	if params.ID == nil || params.RecipientID == nil || params.ExpireString == nil || params.Weight == nil || params.Price == nil || params.Cover == nil {
		fmt.Println("miss required flags")
		return
	}
	err := serv.AcceptFromCourier(model.OrderInput{
		ID:          *params.ID,
		RecipientID: *params.RecipientID,
		Weight:      *params.Weight,
		Price:       *params.Price,
		Cover:       *params.Cover,
		ExpireDate:  *params.ExpireString,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("got new order from courier")
}

func Remove(serv serviceOrder, params parsing.Params) {
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

func Give(serv serviceOrder, params parsing.Params) {
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

func List(serv serviceOrder, params parsing.Params) {
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
		fmt.Printf("%d.\tid: %d\tprice: %d\texpires: %s\n", i+1, order.ID, order.Price, order.ExpireDate.Format("01.02.2006"))
	}
}

func Return(serv serviceOrder, params parsing.Params) {
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

func ListReturn(serv serviceOrder, params parsing.Params) {
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
		fmt.Printf("%d.\tid: %d\trecipient: %d\tprice: %d\texpires: %s\n", startPos+i, order.ID, order.RecipientID, order.Price, order.ExpireDate.Format("01.02.2006"))
	}
}
