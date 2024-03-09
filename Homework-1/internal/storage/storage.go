package storage

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"homework1/pup/internal/model"
	"io"
	"os"
	"time"
)

const storageName = "storage"

type Storage struct {
	storage *os.File
}

// New returns new storage associated with file storageName
func New() (Storage, error) {
	file, err := os.OpenFile(storageName, os.O_CREATE, 0777)
	if err != nil {
		return Storage{}, err
	}
	return Storage{storage: file}, nil
}

// Get adds new order to storage
func (s *Storage) AcceptFromCourier(order model.Order) error {
	if order.ExpireDate.Before(time.Now()) {
		return errors.New("can not get order: trying to get expired order")
	}
	all, err := s.listAll()
	if err != nil {
		return err
	}

	newOrder := OrderDTO{
		ID:          order.ID,
		RecipientID: order.RecipientID,
		ExpireDate:  order.ExpireDate,
	}

	for _, ord := range all {
		if ord.ID == newOrder.ID {
			return errors.New("can not get order: trying to get existing order")
		}
	}
	all = append(all, newOrder)

	return writeBytes(all)
}

// Remove deletes an order from storage
func (s *Storage) Remove(id int64) error {
	all, err := s.listAll()
	if err != nil {
		return err
	}

	for i, ord := range all {
		if ord.ID == id {
			if ord.ExpireDate.Before(time.Now()) && !ord.IsGiven {
				all = append(all[:i], all[i+1:]...)
				return writeBytes(all)
			} else {
				return errors.New("order can not be removed: trying to remove order that is given or not expired")
			}
		}
	}

	return errors.New("order can not be removed: order not found")
}

// Give gives orders to recipient by changing flag IsGiven
func (s *Storage) Give(ids []int64) error {
	all, err := s.listAll()
	if err != nil {
		return err
	}
	var recipient int64
	toModify := make([]int, len(ids))

	for i, id := range ids {
		idx, ok := searchId(all, id)
		if !ok {
			return fmt.Errorf("can not give orders: order %d is not in the storage", id)
		}
		if recipient != 0 && all[idx].RecipientID != recipient {
			return errors.New("can not give orders: orders belong to different recipients")
		}
		if all[idx].IsGiven {
			return fmt.Errorf("can not give orders: order %d is already given", id)
		}
		if all[idx].IsReturned {
			return fmt.Errorf("can not give orders: order %d is already returned by recipient", id)
		}
		if all[idx].ExpireDate.Before(time.Now()) {
			return fmt.Errorf("can not give orders: order %d is expired", id)
		}
		recipient = all[idx].RecipientID
		toModify[i] = idx
	}

	for _, i := range toModify {
		all[i].IsGiven = true
		all[i].GivenTime = time.Now()
	}

	return writeBytes(all)
}

// List returns all recipient's orders
func (s *Storage) List(recipient int64, flag bool) ([]model.Order, error) {
	all, err := s.listAll()
	if err != nil {
		return []model.Order{}, err
	}
	filteredAll := filterOrders(all, func(order OrderDTO) bool {
		return order.RecipientID == recipient && (!flag || !order.IsGiven)
	})
	if len(filteredAll) == 0 {
		return filteredAll, errors.New("can not list orders: orders not found")
	}

	return filteredAll, nil
}

// Return gets order back from recipient by changing flag IsReturned
func (s *Storage) Return(id, recipient int64) error {
	all, err := s.listAll()
	if err != nil {
		return err
	}

	for i, order := range all {
		if order.ID == id {
			if order.RecipientID != recipient {
				return errors.New("order can not be returned: order belongs to different recipient")
			}
			if !order.IsGiven {
				return errors.New("order can not be returned: order is not given yet")
			}
			if order.GivenTime.AddDate(0, 0, 2).Before(time.Now()) {
				return errors.New("order can not be returned: more than 2 days passed")
			}
			all[i].IsGiven = false
			all[i].IsReturned = true
			return writeBytes(all)
		}
	}

	return errors.New("order can not be returned: not found")
}

// ListReturn returns all returned orders in the storage
func (s *Storage) ListReturn() ([]model.Order, error) {
	all, err := s.listAll()
	if err != nil {
		return []model.Order{}, err
	}
	filteredAll := filterOrders(all, func(order OrderDTO) bool {
		return order.IsReturned
	})
	if len(filteredAll) == 0 {
		return filteredAll, errors.New("can not list orders: orders not found")
	}

	return filteredAll, nil
}

// fliterOrders filters slice of orders by filter function
func filterOrders(all []OrderDTO, filter func(OrderDTO) bool) []model.Order {
	filteredAll := make([]model.Order, 0)
	for _, order := range all {
		if filter(order) {
			filteredAll = append(filteredAll, model.Order{
				ID:          order.ID,
				RecipientID: order.RecipientID,
				ExpireDate:  order.ExpireDate,
			})
		}
	}

	return filteredAll
}

// writeBytes writes orders in file in json
func writeBytes(orders []OrderDTO) error {
	rawBytes, err := json.Marshal(orders)
	if err != nil {
		return err
	}

	return os.WriteFile(storageName, rawBytes, 0777)
}

// listAll returns all orders in storage
func (s *Storage) listAll() ([]OrderDTO, error) {
	reader := bufio.NewReader(s.storage)
	rawBytes, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	orders := make([]OrderDTO, 0)
	if len(rawBytes) == 0 {
		return orders, nil
	}

	err = json.Unmarshal(rawBytes, &orders)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

// searchId finds index of order with given id in slice of orders
func searchId(all []OrderDTO, id int64) (int, bool) {
	for i, order := range all {
		if order.ID == id {
			return i, true
		}
	}
	return 0, false
}
