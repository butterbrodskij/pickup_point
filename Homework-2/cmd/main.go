package main

import (
	"fmt"
	"homework2/pup/cmd/command"
	"homework2/pup/cmd/parsing"
	"homework2/pup/internal/service"
	"homework2/pup/internal/storage"
)

func main() {
	var params parsing.Params
	parsing.Parse(&params)

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
		command.Help()
	case "accept":
		command.Accept(serv, params)
	case "remove":
		command.Remove(serv, params)
	case "give":
		command.Give(serv, params)
	case "list":
		command.List(serv, params)
	case "return":
		command.Return(serv, params)
	case "list-return":
		command.ListReturn(serv, params)
	case "pickpoints":
		command.PickPoints(serv)
	default:
		fmt.Println("Unknown command")
	}
}
