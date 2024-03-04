package main

import (
	"flag"
	"fmt"
	"homework1/pup/internal/model"
	"homework1/pup/internal/service"
	"homework1/pup/internal/storage"
)

func help() {
	fmt.Println(`
	usage: go run ./cmd -command=<help|get|remove|give|list|return|list-return> [-id=<order id>] [-recipient=<recipient id>] [-expire=<expire date>] [<args>]

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
		list		 -recipient (optional args: number of orders to list)
		return  	 -id -recipient
		list-return	 args: page number
	
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
	}
}
