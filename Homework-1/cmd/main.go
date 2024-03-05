package main

import (
	"flag"
	"fmt"
	"homework1/pup/internal/model"
	"homework1/pup/internal/service"
	"homework1/pup/internal/storage"
	"strconv"
)

func help() {
	fmt.Println(`
	usage: go run ./cmd -command=<help|get|remove|give|list|return|list-return> [-id=<order id>] [-recipient=<recipient id>] [-expire=<expire date>] [-t=<bool>] [<args>]

	Command desciption:
		help: список доступных команд с кратким описанием
		get: принять заказ от курьера
		remove: вернуть заказ курьеру
		give: выдать заказ клиенту
		list: получить список заказов
		return: принять возврат от клиента
		list-return: получить список возвратов

	Needed flags or arguments for each command:
		help	
		get 		 -id -recipient -expire
		remove  	 -id
		give		 args: order ids to give (example: go run ./cmd -command=give 1 2 3 4)
		list		 -recipient (optional flag -t: boolean value for printing orders located in our point (not already given); optional args: number of orders to list or zero fo all)
		return  	 -id -recipient
		list-return	 args: page number and number of orders per page (default: all pages and 10 orders per page)
	
	Flags requirements:
		-id, -recipient: positive number
		-expire: date in 'dd.mm.yyyy' format (02.01.2006 for 2nd Jan 2006)
	`)
}

func main() {
	command := flag.String("command", "", "name of command")
	id := flag.Int("id", 0, "order id")
	recipient := flag.Int("recipient", 0, "recipient id")
	expireString := flag.String("expire", "", "expire date")
	notGiven := flag.Bool("t", false, "return only not given orders")

	flag.Parse()
	arguments := flag.Args()

	stor, err := storage.New()
	if err != nil {
		fmt.Println("can not connect to storage")
		return
	}
	serv := service.New(&stor)

	switch *command {
	case "":
		fmt.Println("expected a command")
	case "help":
		help()
	case "get":
		if id == nil || recipient == nil || expireString == nil {
			fmt.Println("miss required flags")
			return
		}
		err = serv.Get(model.OrderInput{
			ID:         *id,
			Recipient:  *recipient,
			ExpireDate: *expireString,
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
		err = serv.Give(arguments)
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
		var n int
		k := 10
		if len(arguments) > 0 {
			n, err = strconv.Atoi(arguments[0])
			if err != nil {
				fmt.Println(err)
				return
			}
		}
		if len(arguments) > 1 {
			k, err = strconv.Atoi(arguments[1])
			if err != nil {
				fmt.Println(err)
				return
			}
		}
		arr, err := serv.ListReturn(n, k)
		if err != nil {
			fmt.Println(err)
			return
		}
		startPos := 1
		if n == 0 {
			fmt.Println("all returned not removed orders:")
		} else {
			startPos = k*(n-1) + 1
			fmt.Printf("returned not removed orders from page %d (%d-%d):\n", n, startPos, startPos+len(arr)-1)
		}

		for i, order := range arr {
			fmt.Printf("%d.\tid: %d\trecipient: %d\texpires: %s\n", startPos+i, order.ID, order.Recipient, order.ExpireDate.Format("01.02.2006"))
		}
	}
}
