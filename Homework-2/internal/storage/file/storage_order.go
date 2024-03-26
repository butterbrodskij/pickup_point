package storage

import (
	"encoding/json"
	"errors"
	"os"
	"time"

	"gitlab.ozon.dev/mer_marat/homework/internal/model"
)

type Storage struct {
	storageName string
	content     []OrderDTO
}

// New returns new storage associated with file storageName
func New(storageName string) (Storage, error) {
	_, err := os.Stat(storageName)
	if os.IsNotExist(err) {
		f, err := os.Create(storageName)
		if err != nil {
			return Storage{}, err
		}
		err = f.Close()
		if err != nil {
			return Storage{}, err
		}
		return Storage{
			storageName: storageName,
			content:     []OrderDTO{},
		}, err
	} else if err != nil {
		return Storage{}, err
	}
	content, err := listAll(storageName)
	if err != nil {
		return Storage{}, err
	}
	return Storage{
		storageName: storageName,
		content:     content,
	}, nil
}

// Get returns OrderDTO by its id
func (s *Storage) GetByID(id int64) (OrderDTO, bool) {
	for _, order := range s.content {
		if order.ID == id {
			return order, true
		}
	}
	return OrderDTO{}, false
}

// AcceptFromCourier adds new order to storage
func (s *Storage) AcceptFromCourier(order model.Order) error {
	all := s.content
	newOrder := order2DTO(order)

	for _, ord := range all {
		if ord.ID == newOrder.ID {
			return errors.New("can not get order: trying to get existing order")
		}
	}
	all = append(all, newOrder)

	return s.writeBytes(all)
}

// Remove deletes an order from storage
func (s *Storage) Remove(id int64) error {
	all := s.content

	for i, ord := range all {
		if ord.ID == id {
			all = append(all[:i], all[i+1:]...)
			return s.writeBytes(all)
		}
	}

	return errors.New("order can not be removed: order not found")
}

// Give gives orders to recipient by changing flag IsGiven
func (s *Storage) Give(ids []int64) error {
	all := s.content
	m := make(map[int64]struct{})

	for _, id := range ids {
		m[id] = struct{}{}
	}

	for i, order := range all {
		if _, ok := m[order.ID]; ok {
			all[i].IsGiven = true
			all[i].GivenTime = time.Now()
		}
	}

	return s.writeBytes(all)
}

// List returns all recipient's orders
func (s *Storage) ListAll(recipient int64) ([]model.Order, error) {
	all := s.content
	filteredAll := filterOrders(all, func(order OrderDTO) bool {
		return order.RecipientID == recipient
	})
	if len(filteredAll) == 0 {
		return filteredAll, errors.New("can not list orders: orders not found")
	}

	return filteredAll, nil
}

// List returns all recipient's not given orders
func (s *Storage) ListNotGiven(recipient int64) ([]model.Order, error) {
	all := s.content
	filteredAll := filterOrders(all, func(order OrderDTO) bool {
		return order.RecipientID == recipient && !order.IsGiven
	})
	if len(filteredAll) == 0 {
		return filteredAll, errors.New("can not list orders: orders not found")
	}

	return filteredAll, nil
}

// Return gets order back from recipient by changing flag IsReturned
func (s *Storage) Return(id int64) error {
	all := s.content

	for i, order := range all {
		if order.ID == id {
			all[i].IsGiven = false
			all[i].IsReturned = true
			return s.writeBytes(all)
		}
	}

	return errors.New("order can not be returned: order not found")
}

// ListReturn returns all returned orders in the storage
func (s *Storage) ListReturn() ([]model.Order, error) {
	all := s.content
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
			filteredAll = append(filteredAll, dto2Order(order))
		}
	}

	return filteredAll
}

// writeBytes writes orders in file in json
func (s *Storage) writeBytes(orders []OrderDTO) error {
	s.content = orders
	rawBytes, err := json.Marshal(orders)
	if err != nil {
		return err
	}

	return os.WriteFile(s.storageName, rawBytes, 0777)
}

// listAll returns all orders in storage
func listAll(storageName string) ([]OrderDTO, error) {
	rawBytes, err := os.ReadFile(storageName)
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

func order2DTO(order model.Order) OrderDTO {
	return OrderDTO{
		ID:          order.ID,
		RecipientID: order.RecipientID,
		Weight:      order.Weight,
		Price:       order.Price,
		Cover:       order.Cover,
		ExpireDate:  order.ExpireDate,
	}
}

func dto2Order(order OrderDTO) model.Order {
	return model.Order{
		ID:          order.ID,
		RecipientID: order.RecipientID,
		Weight:      order.Weight,
		Price:       order.Price,
		Cover:       order.Cover,
		ExpireDate:  order.ExpireDate,
	}
}
