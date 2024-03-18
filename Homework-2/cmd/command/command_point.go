package command

import (
	"bufio"
	"context"
	"fmt"
	"homework2/pup/internal/model"
	"homework2/pup/internal/service/pickpoint"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

const chanSize = 10

// Implementation of command pickpoints
func PickPoints(serv pickpoint.Service) {
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

	go WritePoints(serv, ctx, writeChan, logWriteChan, &wg)
	go Reader(serv, ctx, readChan, logReadChan, &wg)
	go LogPoints(serv, ctx, logWriteChan, logReadChan, &wg)

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

// Reader makes pool of readers
func Reader(s pickpoint.Service, ctx context.Context, readChan <-chan int64, logChan chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	var wgReader sync.WaitGroup
	wgReader.Add(chanSize)
	for i := 1; i <= chanSize; i++ {
		serial := i
		go ReadPoints(s, ctx, readChan, logChan, &wgReader, serial)
	}
	wgReader.Wait()
	close(logChan)
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

// WritePoints writes pick-up points information in storage from channel
func WritePoints(s pickpoint.Service, ctx context.Context, writeChan <-chan model.PickPoint, logChan chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	var status string
	for {
		select {
		case <-ctx.Done():
			message := "writer: context is canceled"
			logChan <- message
			close(logChan)
			return
		case point := <-writeChan:
			message := fmt.Sprintf("writer: trying to write new pick-up point %v", point)
			logChan <- message
			err := s.Write(point)
			if err != nil {
				status = fmt.Sprintf("writer: error while adding point %d: %s", point.ID, err.Error())
			} else {
				status = fmt.Sprintf("writer: point %d added successfully", point.ID)
			}
			logChan <- status
		}
	}
}

// WritePoints sends pick-up points information to logger from storage by getting id from channel
func ReadPoints(s pickpoint.Service, ctx context.Context, readChan <-chan int64, logChan chan<- string, wg *sync.WaitGroup, serial int) {
	defer wg.Done()
	var status string
	for {
		select {
		case <-ctx.Done():
			message := fmt.Sprintf("reader %d: context is canceled", serial)
			logChan <- message
			return
		case id := <-readChan:
			message := fmt.Sprintf("reader %d: trying to find info about pick-up point with id %d", serial, id)
			logChan <- message
			point, err := s.Get(id)
			if err != nil {
				status = fmt.Sprintf("reader %d: error while getting point %d: %s", serial, id, err)
			} else {
				status = fmt.Sprintf("reader %d: found pick-up point:\n\tid: %d\tname: %s\taddress: %s\tcontacts: %s", serial, point.ID, point.Name, point.Address, point.Contact)
			}
			logChan <- status
		}
	}
}

// LogPoints prints all logs from writer and reader
func LogPoints(s pickpoint.Service, ctx context.Context, logWriteChan, logReadChan <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			for s := range logReadChan {
				fmt.Println(s)
			}
			for s := range logWriteChan {
				fmt.Println(s)
			}
			fmt.Println("logger: context is canceled")
			return
		case s := <-logWriteChan:
			fmt.Println(s)
		case s := <-logReadChan:
			fmt.Println(s)
		}
	}
}
