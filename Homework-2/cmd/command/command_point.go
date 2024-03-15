package command

import (
	"bufio"
	"context"
	"fmt"
	"homework2/pup/internal/model"
	"homework2/pup/internal/service"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

const chanSize = 10

// Implementation of command pickpoints
func PickPoints(serv service.Service) {
	var (
		line, com string
		id        int64
		point     model.PickPoint
		wg        sync.WaitGroup
	)
	ctx, cancel := context.WithCancel(context.Background())
	wg.Add(3)
	writeChan := make(chan model.PickPoint, chanSize)
	readChan := make(chan int64, chanSize)
	logReadChan := make(chan string, chanSize)
	logWriteChan := make(chan string, chanSize)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go serv.WritePoints(ctx, writeChan, logWriteChan, &wg)
	go serv.ReadPoints(ctx, readChan, logReadChan, &wg)
	go serv.LogPoints(ctx, logWriteChan, logReadChan, &wg)

	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			line = scanner.Text()
			fmt.Sscanf(line, "%s", &com)
			switch com {
			case "help":
				HelpPickPoints()
			case "exit":
				cancel()
				return
			case "write":
				_, err := fmt.Sscanf(line, "write %d %s %s %s", &point.ID, &point.Name, &point.Address, &point.Contact)
				if err != nil {
					fmt.Println(err)
					continue
				}
				writeChan <- point
			case "read":
				_, err := fmt.Sscanf(line, "read %d", &id)
				if err != nil {
					fmt.Println(err)
					continue
				}
				readChan <- id
			default:
				fmt.Println("Unknown command")
			}
		}
	}()

	for {
		select {
		case sig := <-sigChan:
			fmt.Println("\nsignal caught:", sig)
			cancel()
			wg.Wait()
			return
		case <-ctx.Done():
			wg.Wait()
			return
		}
	}
}

// HelpPickPoints prints usage guide for pickpoints
func HelpPickPoints() {
	fmt.Println(`
	interactive mode for command pickpoints usage guide:

	Command desciption:
		help: список доступных команд с кратким описанием
		write: добавить информацию о ПВЗ
		read: считать информацию о ПВЗ
		exit: завершение работы

	Needed arguments for each command:
		help	
		write 		 id(int)	name(string)	address(string)	   contact(string)
		read	  	 id(int)
		exit
	
	Examples:
		write 10 Chertanovo Chertanovskaya-Street-10 +7(999)888-77-66
		read 10
	`)
}
