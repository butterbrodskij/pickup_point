package order

import (
	"errors"
	"fmt"
	"time"

	"gitlab.ozon.dev/mer_marat/homework/internal/model"
	"gitlab.ozon.dev/mer_marat/homework/internal/service/cover"
	storage "gitlab.ozon.dev/mer_marat/homework/internal/storage/file"
)

type storageInterface interface {
	AcceptFromCourier(model.Order) error
	Remove(int64) error
	Give([]int64) error
	ListAll(int64) ([]model.Order, error)
	ListNotGiven(int64) ([]model.Order, error)
	Return(int64) error
	ListReturn() ([]model.Order, error)
	GetByID(int64) (storage.OrderDTO, bool)
}

type service struct {
	s storageInterface
}

// Input2Order converts OrderInput to Order and checks validity of fields
func input2Order(input model.OrderInput) (model.Order, error) {
	if input.ID <= 0 {
		return model.Order{}, errors.New("id should be positive")
	}
	if input.RecipientID <= 0 {
		return model.Order{}, errors.New("recipient id should be positive")
	}
	if input.Weight <= 0 {
		return model.Order{}, errors.New("order weight should be positive")
	}
	if input.Price <= 0 {
		return model.Order{}, errors.New("order price should be positive")
	}

	t, err := time.Parse("2.1.2006", input.ExpireDate)
	if err != nil {
		return model.Order{}, errors.New("wrong date format")
	}

	return model.Order{
		ID:          input.ID,
		RecipientID: input.RecipientID,
		Weight:      input.Weight,
		Price:       input.Price,
		Cover:       input.Cover,
		ExpireDate:  t,
	}, nil
}

// New returns type Service associated with storage
func NewService(stor storageInterface) service {
	return service{s: stor}
}

// Get checks validity of given data and adds new order to storage
func (s service) AcceptFromCourier(input model.OrderInput) error {
	order, err := input2Order(input)
	if err != nil {
		return err
	}
	if order.ExpireDate.Before(time.Now()) {
		return errors.New("can not get order: trying to get expired order")
	}
	covered, err := cover.CoveredOrder(&order)
	if err != nil {
		return err
	}
	if ok := covered.OrderRequirements(); !ok {
		return errors.New("order does not meet cover requirements")
	}
	return s.s.AcceptFromCourier(*covered.OrderChanges())
}

// Remove checks validity of given id and deletes an order from storage
func (s service) Remove(id int64) error {
	if id <= 0 {
		return errors.New("id should be positive")
	}
	if order, ok := s.s.GetByID(id); ok && order.ExpireDate.After(time.Now()) || order.IsGiven {
		return errors.New("order can not be removed: trying to remove order that is given or not expired")
	}
	return s.s.Remove(id)
}

// Give checks validity of given ids and gives orders to recipient
func (s service) Give(ids []int64) error {
	var recipient int64

	for _, id := range ids {
		order, ok := s.s.GetByID(id)
		if !ok {
			return fmt.Errorf("can not give orders: order %d is not in the storage", id)
		}
		if recipient != 0 && order.RecipientID != recipient {
			return errors.New("can not give orders: orders belong to different recipients")
		}
		if order.IsGiven {
			return fmt.Errorf("can not give orders: order %d is already given", id)
		}
		if order.IsReturned {
			return fmt.Errorf("can not give orders: order %d is already returned by recipient", id)
		}
		if order.ExpireDate.Before(time.Now()) {
			return fmt.Errorf("can not give orders: order %d is expired", id)
		}
		if recipient == 0 {
			recipient = order.RecipientID
		}
	}

	return s.s.Give(ids)
}

// List checks validity of given recipient id and n and returns slice of all his orders (last n)
func (s service) List(recipient int64, n int, onlyNotGivenOrders bool) ([]model.Order, error) {
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
	if err != nil {
		return []model.Order{}, err
	}
	if n == 0 || len(all) <= n {
		return all, nil
	}

	return all[len(all)-n:], err
}

// Return checks validity of given order id and recipient id and gets order back from recipient
func (s service) Return(id, recipient int64) error {
	if id <= 0 {
		return errors.New("id should be positive")
	}
	if recipient <= 0 {
		return errors.New("recipient id should be positive")
	}
	order, ok := s.s.GetByID(id)
	if !ok {
		return errors.New("order can not be returned: order not found")
	}
	if order.RecipientID != recipient {
		return errors.New("order can not be returned: order belongs to different recipient")
	}
	if order.IsReturned {
		return errors.New("order can not be returned: order is already returned")
	}
	if !order.IsGiven {
		return errors.New("order can not be returned: order is not given yet")
	}
	if order.GivenTime.AddDate(0, 0, 2).Before(time.Now()) {
		return errors.New("order can not be returned: more than 2 days passed")
	}
	return s.s.Return(id)
}

// ListReturn checks validity of given args and returns k returned orders on nth page
func (s service) ListReturn(pageNum, ordersPerPage int) ([]model.Order, error) {
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
