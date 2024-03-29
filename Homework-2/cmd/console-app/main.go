package main

import (
	"fmt"

	"gitlab.ozon.dev/mer_marat/homework/cmd/console-app/command"
	"gitlab.ozon.dev/mer_marat/homework/cmd/console-app/parsing"
	"gitlab.ozon.dev/mer_marat/homework/internal/service/cover"
	"gitlab.ozon.dev/mer_marat/homework/internal/service/order"
	"gitlab.ozon.dev/mer_marat/homework/internal/service/pickpoint"
	storage "gitlab.ozon.dev/mer_marat/homework/internal/storage/file"
)

func main() {
	var params parsing.Params
	parsing.Parse(&params)

	storOrders, err := storage.NewOrders("storage_orders.json")
	if err != nil {
		fmt.Printf("can not connect to storage: %s\n", err)
		return
	}
	storPoints, err := storage.NewPoints("storage_points.json")
	if err != nil {
		fmt.Printf("can not connect to storage: %s\n", err)
		return
	}
	servOrders := order.NewService(&storOrders, cover.NewService())
	servPoints := pickpoint.NewService(&storPoints)

	switch *params.Command {
	case "":
		fmt.Println("expected a command")
	case "help":
		command.Help()
	case "accept":
		command.Accept(servOrders, params)
	case "remove":
		command.Remove(servOrders, params)
	case "give":
		command.Give(servOrders, params)
	case "list":
		command.List(servOrders, params)
	case "return":
		command.Return(servOrders, params)
	case "list-return":
		command.ListReturn(servOrders, params)
	case "pickpoints":
		command.PickPoints(servPoints)
	default:
		fmt.Println("Unknown command")
	}
}
