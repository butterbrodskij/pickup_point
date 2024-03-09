package main

import "fmt"

// help prints usage guide
func help() {
	fmt.Println(`
	usage: go run ./cmd -command=<help|accept|remove|give|list|return|list-return> [-id=<order id>] [-recipient=<recipient id>] [-expire=<expire date>] [-t=<bool>] [<args>]

	Command desciption:
		help: список доступных команд с кратким описанием
		accept: принять заказ от курьера
		remove: вернуть заказ курьеру
		give: выдать заказ клиенту
		list: получить список заказов
		return: принять возврат от клиента
		list-return: получить список возвратов

	Needed flags or arguments for each command:
		help	
		accept 		 -id -recipient -expire
		remove  	 -id
		give		 args: order ids to give (example: go run ./cmd -command=give 1 2 3 4)
		list		 -recipient (optional flag -t: boolean value for printing orders located in our point (not already given); optional args: number of orders to list or zero fo all)
		return  	 -id -recipient
		list-return	 args: page number and number of orders per page (default: all pages and 10 orders per page) (example: "-command=list-return 2 5" prints 2nd page of returned orders grouped by 5 orders in each page)
	
	Flags requirements:
		-id, -recipient: positive number
		-expire: date in 'dd.mm.yyyy' format (02.01.2006 for 2nd Jan 2006)
	`)
}
