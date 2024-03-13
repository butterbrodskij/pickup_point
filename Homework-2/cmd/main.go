package main

import (
	"fmt"
	"homework2/pup/cmd/flags"
	"homework2/pup/internal/model"
	"homework2/pup/internal/service"
	"homework2/pup/internal/storage"
	"strconv"
)

func main() {
	var params flags.Params
	flags.Parse(&params)

	stor, err := storage.New("storage.json")
	if err != nil {
		fmt.Printf("can not connect to storage: %s\n", err)
		return
	}
	serv := service.New(&stor)

	switch *params.Command {
	case "":
		fmt.Println("expected a command")
	case "help":
		serv.Help()
	case "accept":
		if params.ID == nil || params.RecipientID == nil || params.ExpireString == nil {
			fmt.Println("miss required flags")
			return
		}
		err = serv.AcceptFromCourier(model.OrderInput{
			ID:          *params.ID,
			RecipientID: *params.RecipientID,
			ExpireDate:  *params.ExpireString,
		})
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("got new order from courier")
	case "remove":
		if params.ID == nil || params.RecipientID == nil || params.ExpireString == nil {
			fmt.Println("miss required flags")
			return
		}
		err = serv.Remove(*params.ID)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("removed order %d from our pick-up point\n", *params.ID)
	case "give":
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
		err = serv.Give(ids)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("all orders have been given to its recipient")
	case "list":
		if params.RecipientID == nil {
			fmt.Println("miss required flags")
			return
		}
		var n int
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
	case "return":
		if params.ID == nil || params.RecipientID == nil {
			fmt.Println("miss required flags")
			return
		}
		err = serv.Return(*params.ID, *params.RecipientID)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("order %d is returned successfully\n", *params.ID)
	case "list-return":
		var pageNum int
		ordersPerPage := 10
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
}
