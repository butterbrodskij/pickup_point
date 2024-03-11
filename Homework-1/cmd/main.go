package main

import (
	"flag"
	"fmt"
	"homework1/pup/internal/model"
	"homework1/pup/internal/service"
	"homework1/pup/internal/storage"
	"strconv"
)

func main() {
	command := flag.String("command", "", "name of command")
	id := flag.Int64("id", 0, "order id")
	recipient := flag.Int64("recipient", 0, "recipient id")
	expireString := flag.String("expire", "", "expire date")
	notGiven := flag.Bool("t", false, "return only not given orders")

	flag.Parse()
	arguments := flag.Args()

	stor, err := storage.New("storage.json")
	if err != nil {
		fmt.Printf("can not connect to storage: %s\n", err)
		return
	}
	serv := service.New(&stor)

	switch *command {
	case "":
		fmt.Println("expected a command")
	case "help":
		serv.Help()
	case "accept":
		if id == nil || recipient == nil || expireString == nil {
			fmt.Println("miss required flags")
			return
		}
		err = serv.AcceptFromCourier(model.OrderInput{
			ID:          *id,
			RecipientID: *recipient,
			ExpireDate:  *expireString,
		})
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("got new order from courier")
	case "remove":
		if id == nil || recipient == nil || expireString == nil {
			fmt.Println("miss required flags")
			return
		}
		err = serv.Remove(*id)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("removed order %d from our pick-up point\n", *id)
	case "give":
		if len(arguments) == 0 {
			fmt.Println("expected at least one argument as order id")
			return
		}
		ids := make([]int64, len(arguments))
		for i, s := range arguments {
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
		if recipient == nil {
			fmt.Println("miss required flags")
			return
		}
		var n int
		if len(arguments) > 0 {
			n, err = strconv.Atoi(arguments[0])
			if err != nil {
				fmt.Println(err)
				return
			}
		}
		foundArr, err := serv.List(*recipient, n, *notGiven)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("found %d orders:\n", len(foundArr))

		for i, order := range foundArr {
			fmt.Printf("%d.\tid: %d\texpires: %s\n", i+1, order.ID, order.ExpireDate.Format("01.02.2006"))
		}
	case "return":
		if id == nil || recipient == nil {
			fmt.Println("miss required flags")
			return
		}
		err = serv.Return(*id, *recipient)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("order %d is returned successfully\n", *id)
	case "list-return":
		var pageNum int
		ordersPerPage := 10
		if len(arguments) > 0 {
			pageNum, err = strconv.Atoi(arguments[0])
			if err != nil {
				fmt.Println(err)
				return
			}
		}
		if len(arguments) > 1 {
			ordersPerPage, err = strconv.Atoi(arguments[1])
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
