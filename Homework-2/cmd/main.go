package main

import (
	"fmt"
	"homework2/pup/cmd/command"
	"homework2/pup/cmd/parsing"
	"homework2/pup/internal/service/order"
	"homework2/pup/internal/service/pickpoint"
	"homework2/pup/internal/storage"
)

func main() {
	var params parsing.Params
	parsing.Parse(&params)

	storOrders, err := storage.New("storage.json")
	if err != nil {
		fmt.Printf("can not connect to storage: %s\n", err)
		return
	}
	storPoints, err := storage.NewPoints("storage_points.json")
	if err != nil {
		fmt.Printf("can not connect to storage: %s\n", err)
		return
	}
	servOrders := order.New(&storOrders)
	servPoints := pickpoint.New(&storPoints)

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
