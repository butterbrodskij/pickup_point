package service

import (
	"errors"
	"homework1/pup/internal/model"
	"time"
)

type storage interface {
	AcceptFromCourier(model.Order) error
	Remove(int64) error
	Give([]int64) error
	ListAll(int64) ([]model.Order, error)
	ListNotGiven(int64) ([]model.Order, error)
	Return(int64, int64) error
	ListReturn() ([]model.Order, error)
}

type Service struct {
	s storage
}

// Input2Order converts OrderInput to Order and checks validity of fields
func Input2Order(input model.OrderInput) (model.Order, error) {
	if input.ID <= 0 {
		return model.Order{}, errors.New("id should be positive")
	}

	if input.RecipientID <= 0 {
		return model.Order{}, errors.New("recipient id should be positive")
	}

	t, err := time.Parse("2.1.2006", input.ExpireDate)
	if err != nil {
		return model.Order{}, errors.New("wrong date format")
	}

	return model.Order{
		ID:          input.ID,
		RecipientID: input.RecipientID,
		ExpireDate:  t,
	}, nil
}

// New returns type Service associated with storage
func New(stor storage) Service {
	return Service{s: stor}
}

// Get checks validity of given data and adds new order to storage
func (s Service) AcceptFromCourier(input model.OrderInput) error {
	order, err := Input2Order(input)
	if err != nil {
		return err
	}
	if order.ExpireDate.Before(time.Now()) {
		return errors.New("can not get order: trying to get expired order")
	}
	return s.s.AcceptFromCourier(order)
}

// Remove checks validity of given id and deletes an order from storage
func (s Service) Remove(id int64) error {
	if id <= 0 {
		return errors.New("id should be positive")
	}
	return s.s.Remove(id)
}

// Give checks validity of given ids and gives orders to recipient
func (s Service) Give(ids []int64) error {
	return s.s.Give(ids)
}

// List checks validity of given recipient id and n and returns slice of all his orders (last n)
func (s Service) List(recipient int64, n int, onlyNotGivenOrders bool) ([]model.Order, error) {
	if recipient <= 0 {
		return []model.Order{}, errors.New("recipient id should be positive")
	}
	if n < 0 {
		return []model.Order{}, errors.New("n should not be negative")
	}
	var (
		all []model.Order
		err error
	)
	if !onlyNotGivenOrders {
		all, err = s.s.ListAll(recipient)
	} else {
		all, err = s.s.ListNotGiven(recipient)
	}
	if err != nil || n == 0 || len(all) <= n {
		return all, err
	}

	return all[len(all)-n:], err
}

// Return checks validity of given order id and recipient id and gets order back from recipient
func (s Service) Return(id, recipient int64) error {
	if id <= 0 {
		return errors.New("id should be positive")
	}
	if recipient <= 0 {
		return errors.New("recipient id should be positive")
	}
	return s.s.Return(id, recipient)
}

// ListReturn checks validity of given args and returns k returned orders on nth page
func (s Service) ListReturn(pageNum, ordersPerPage int) ([]model.Order, error) {
	if pageNum < 0 {
		return []model.Order{}, errors.New("pageNum should not be negative")
	}
	if ordersPerPage <= 0 {
		return []model.Order{}, errors.New("ordersPerPage should be positive")
	}
	all, err := s.s.ListReturn()
	if err != nil || pageNum == 0 {
		return all, err
	}
	firstPos := (pageNum - 1) * ordersPerPage
	if len(all) == 0 || len(all) <= firstPos {
		return all, errors.New("empty list")
	}
	newLen := ordersPerPage
	if len(all) < pageNum*ordersPerPage {
		newLen = len(all) % ordersPerPage
	}
	return all[firstPos : firstPos+newLen], nil
}
