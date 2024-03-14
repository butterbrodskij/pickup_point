package service

import (
	"context"
	"fmt"
	"homework2/pup/internal/model"
	"sync"
)

type storagePointsInterface interface {
	//Write(model.PickPoint) error
	//Get(int64) (model.PickPoint, bool)
}

func (s Service) WritePoints(ctx context.Context, writeChan <-chan model.PickPoint, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			fmt.Println("writer: context is canceled")
			return
		case point := <-writeChan:
			fmt.Println("writer:", point)
		}
	}
}

func (s Service) ReadPoints(ctx context.Context, readChan <-chan int64, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			fmt.Println("reader: context is canceled")
			return
		case id := <-readChan:
			fmt.Println("reader:", id)
		}
	}
}
