package service

import (
	"context"
	"fmt"
	"homework2/pup/internal/model"
	"sync"
)

type storagePointsInterface interface {
	Write(model.PickPoint) error
	Get(int64) (model.PickPoint, bool)
}

// WritePoints writes pick-up points information in storage from channel
func (s Service) WritePoints(ctx context.Context, writeChan <-chan model.PickPoint, logChan chan<- string, wg *sync.WaitGroup) {
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
			message := fmt.Sprint("writer: trying to write new pick-up point ", point)
			logChan <- message
			err := s.sPoints.Write(point)
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
func (s Service) ReadPoints(ctx context.Context, readChan <-chan int64, logChan chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	var status string
	for {
		select {
		case <-ctx.Done():
			message := "reader: context is canceled"
			logChan <- message
			close(logChan)
			return
		case id := <-readChan:
			message := fmt.Sprint("reader: trying to find info about pick-up point with id ", id)
			logChan <- message
			point, ok := s.sPoints.Get(id)
			if !ok {
				status = fmt.Sprintf("reader: point %d not found", id)
			} else {
				status = fmt.Sprintf("reader: found pick-up point:\n\tid: %d\tname: %s\taddress: %s\tcontacts: %s", point.ID, point.Name, point.Address, point.Contact)
			}
			logChan <- status
		}
	}
}

// LogPoints prints all logs from writer and reader
func (s Service) LogPoints(ctx context.Context, logWriteChan, logReadChan <-chan string, wg *sync.WaitGroup) {
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
