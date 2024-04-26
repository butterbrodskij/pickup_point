package main

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"gitlab.ozon.dev/mer_marat/homework/cmd/console-app/command"
	"gitlab.ozon.dev/mer_marat/homework/cmd/console-app/parsing"
	inmemorycache "gitlab.ozon.dev/mer_marat/homework/internal/pkg/in_memory_cache"
	"gitlab.ozon.dev/mer_marat/homework/internal/pkg/transactor"
	"gitlab.ozon.dev/mer_marat/homework/internal/service/cover"
	"gitlab.ozon.dev/mer_marat/homework/internal/service/order"
	"gitlab.ozon.dev/mer_marat/homework/internal/service/pickpoint"
	storage "gitlab.ozon.dev/mer_marat/homework/internal/storage/file"
)

var (
	reg = prometheus.NewRegistry()

	pickpointCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "pickpoint_cli",
		Help: "Number of requests handled",
	})
)

func init() {
	reg.MustRegister(pickpointCounter)
}

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
	cache := inmemorycache.NewInMemoryCache()
	servPoints := pickpoint.NewService(&storPoints, cache, transactor.NewDummyTransactor(), pickpointCounter)

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
