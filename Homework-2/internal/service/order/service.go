package order

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"gitlab.ozon.dev/mer_marat/homework/internal/model"
	storage "gitlab.ozon.dev/mer_marat/homework/internal/storage/file"
)

type coverService interface {
	ValidateOrder(model.Order) error
	GetPackagingPrice(model.Order) int64
}

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
	s                  storageInterface
	cov                coverService
	givenOrdersCounter prometheus.Gauge
	failedCounter      prometheus.Counter
}

func input2Order(input model.OrderInput) (model.Order, error) {
	if input.ID <= 0 {
		return model.Order{}, errors.New("id should be positive")
	}
	if input.RecipientID <= 0 {
		return model.Order{}, errors.New("recipient id should be positive")
	}
	if input.WeightGrams <= 0 {
		return model.Order{}, errors.New("order weight should be positive")
	}
	if input.PriceKopecks <= 0 {
		return model.Order{}, errors.New("order price should be positive")
	}

	t, err := time.Parse("2.1.2006", input.ExpireDate)
	if err != nil {
		return model.Order{}, errors.New("wrong date format")
	}

	return model.Order{
		ID:           input.ID,
		RecipientID:  input.RecipientID,
		WeightGrams:  input.WeightGrams,
		PriceKopecks: input.PriceKopecks,
		Cover:        input.Cover,
		ExpireDate:   t,
	}, nil
}

// New returns type Service associated with storage
func NewService(stor storageInterface, cov coverService) service {
	return service{s: stor, cov: cov}
}

func (s *service) AddGivenOrdersGauge(metrics prometheus.Gauge) {
	s.givenOrdersCounter = metrics
}

func (s *service) AddFailedRequestsCounter(counter prometheus.Counter) {
	s.failedCounter = counter
}

func (s *service) incFailedOrdersCounter() {
	if s.failedCounter != nil {
		s.failedCounter.Inc()
	}
}

// Get checks validity of given data and adds new order to storage
func (s service) AcceptFromCourier(ctx context.Context, input model.OrderInput) error {
	order, err := input2Order(input)
	if err != nil {
		s.incFailedOrdersCounter()
		return err
	}
	if order.ExpireDate.Before(time.Now()) {
		s.incFailedOrdersCounter()
		return errors.New("can not get order: trying to get expired order")
	}
	if err = s.cov.ValidateOrder(order); err != nil {
		s.incFailedOrdersCounter()
		return err
	}
	order.PriceKopecks += s.cov.GetPackagingPrice(order)
	if err = s.s.AcceptFromCourier(order); err != nil {
		s.incFailedOrdersCounter()
		return err
	}
	return nil
}

// Remove checks validity of given id and deletes an order from storage
func (s service) Remove(ctx context.Context, id int64) error {
	if id <= 0 {
		s.incFailedOrdersCounter()
		return errors.New("id should be positive")
	}
	if order, ok := s.s.GetByID(id); ok && order.ExpireDate.After(time.Now()) || order.IsGiven {
		s.incFailedOrdersCounter()
		return errors.New("order can not be removed: trying to remove order that is given or not expired")
	}
	if err := s.s.Remove(id); err != nil {
		s.incFailedOrdersCounter()
		return err
	}
	return nil
}

// Give checks validity of given ids and gives orders to recipient
func (s service) Give(ctx context.Context, ids []int64) error {
	var recipient int64

	for _, id := range ids {
		order, ok := s.s.GetByID(id)
		switch {
		case !ok:
			s.incFailedOrdersCounter()
			return fmt.Errorf("can not give orders: order %d is not in the storage", id)
		case recipient != 0 && order.RecipientID != recipient:
			s.incFailedOrdersCounter()
			return errors.New("can not give orders: orders belong to different recipients")
		case order.IsGiven:
			s.incFailedOrdersCounter()
			return fmt.Errorf("can not give orders: order %d is already given", id)
		case order.IsReturned:
			s.incFailedOrdersCounter()
			return fmt.Errorf("can not give orders: order %d is already returned by recipient", id)
		case order.ExpireDate.Before(time.Now()):
			s.incFailedOrdersCounter()
			return fmt.Errorf("can not give orders: order %d is expired", id)
		case recipient == 0:
			recipient = order.RecipientID
		}
	}
	err := s.s.Give(ids)
	if err != nil {
		s.incFailedOrdersCounter()
		return err
	}
	if s.givenOrdersCounter != nil {
		s.givenOrdersCounter.Add(float64(len(ids)))
	}

	return nil
}

// List checks validity of given recipient id and n and returns slice of all his orders (last n)
func (s service) List(ctx context.Context, recipient int64, n int, onlyNotGivenOrders bool) ([]model.Order, error) {
	if recipient <= 0 {
		s.incFailedOrdersCounter()
		return nil, errors.New("recipient id should be positive")
	}
	if n < 0 {
		s.incFailedOrdersCounter()
		return nil, errors.New("n should not be negative")
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
		s.incFailedOrdersCounter()
		return nil, err
	}
	if n == 0 || len(all) <= n {
		return all, nil
	}

	return all[len(all)-n:], nil
}

// Return checks validity of given order id and recipient id and gets order back from recipient
func (s service) Return(ctx context.Context, id, recipient int64) error {
	if id <= 0 {
		s.incFailedOrdersCounter()
		return errors.New("id should be positive")
	}
	if recipient <= 0 {
		s.incFailedOrdersCounter()
		return errors.New("recipient id should be positive")
	}
	order, ok := s.s.GetByID(id)
	switch {
	case !ok:
		s.incFailedOrdersCounter()
		return errors.New("order can not be returned: order not found")
	case order.RecipientID != recipient:
		s.incFailedOrdersCounter()
		return errors.New("order can not be returned: order belongs to different recipient")
	case order.IsReturned:
		s.incFailedOrdersCounter()
		return errors.New("order can not be returned: order is already returned")
	case !order.IsGiven:
		s.incFailedOrdersCounter()
		return errors.New("order can not be returned: order is not given yet")
	case order.GivenTime.AddDate(0, 0, 2).Before(time.Now()):
		s.incFailedOrdersCounter()
		return errors.New("order can not be returned: more than 2 days passed")
	}
	err := s.s.Return(id)
	if err != nil {
		s.incFailedOrdersCounter()
		return err
	}
	if s.givenOrdersCounter != nil {
		s.givenOrdersCounter.Sub(1)
	}

	return nil
}

// ListReturn checks validity of given args and returns k returned orders on nth page
func (s service) ListReturn(ctx context.Context, pageNum, ordersPerPage int) ([]model.Order, error) {
	if pageNum < 0 {
		s.incFailedOrdersCounter()
		return nil, errors.New("pageNum should not be negative")
	}
	if ordersPerPage <= 0 {
		s.incFailedOrdersCounter()
		return nil, errors.New("ordersPerPage should be positive")
	}
	all, err := s.s.ListReturn()
	if err != nil {
		s.incFailedOrdersCounter()
		return nil, err
	}
	if pageNum == 0 {
		return all, nil
	}
	firstPos := (pageNum - 1) * ordersPerPage
	if len(all) == 0 || len(all) <= firstPos {
		s.incFailedOrdersCounter()
		return nil, errors.New("empty list")
	}
	newLen := ordersPerPage
	if len(all) < pageNum*ordersPerPage {
		newLen = len(all) % ordersPerPage
	}
	return all[firstPos : firstPos+newLen], nil
}
