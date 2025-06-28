package order

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gitlab.ozon.dev/mer_marat/homework/internal/metrics"
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

type gauge interface {
	SuccessGaugeAdd(err error, add float64)
	SuccessGaugeDec(err error)
}

type counter interface {
	FailedCounterInc(err error)
}

type service struct {
	s                storageInterface
	cov              coverService
	givenOrdersGauge gauge
	failedCounter    counter
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
	return service{
		s:                stor,
		cov:              cov,
		givenOrdersGauge: &metrics.UnImplementedGauge{},
		failedCounter:    &metrics.UnImplementedCounter{},
	}
}

func (s *service) AddGivenOrdersGauge(gauge gauge) {
	s.givenOrdersGauge = gauge
}

func (s *service) AddFailedRequestsCounter(counter counter) {
	s.failedCounter = counter
}

// Get checks validity of given data and adds new order to storage
func (s service) AcceptFromCourier(ctx context.Context, input model.OrderInput) (err error) {
	defer func() {
		s.failedCounter.FailedCounterInc(err)
	}()
	order, err := input2Order(input)
	if err != nil {
		return err
	}
	if order.ExpireDate.Before(time.Now()) {
		return errors.New("can not get order: trying to get expired order")
	}
	if err = s.cov.ValidateOrder(order); err != nil {
		return err
	}
	order.PriceKopecks += s.cov.GetPackagingPrice(order)
	if err = s.s.AcceptFromCourier(order); err != nil {
		return err
	}
	return nil
}

// Remove checks validity of given id and deletes an order from storage
func (s service) Remove(ctx context.Context, id int64) (err error) {
	defer func() {
		s.failedCounter.FailedCounterInc(err)
	}()
	if id <= 0 {
		return errors.New("id should be positive")
	}
	if order, ok := s.s.GetByID(id); ok && order.ExpireDate.After(time.Now()) || order.IsGiven {
		return errors.New("order can not be removed: trying to remove order that is given or not expired")
	}
	if err := s.s.Remove(id); err != nil {
		return err
	}
	return nil
}

// Give checks validity of given ids and gives orders to recipient
func (s service) Give(ctx context.Context, ids []int64) (err error) {
	defer func() {
		s.failedCounter.FailedCounterInc(err)
	}()
	defer func() {
		defer s.givenOrdersGauge.SuccessGaugeAdd(err, float64(len(ids)))
	}()
	var recipient int64

	for _, id := range ids {
		order, ok := s.s.GetByID(id)
		switch {
		case !ok:
			return fmt.Errorf("can not give orders: order %d is not in the storage", id)
		case recipient != 0 && order.RecipientID != recipient:
			return errors.New("can not give orders: orders belong to different recipients")
		case order.IsGiven:
			return fmt.Errorf("can not give orders: order %d is already given", id)
		case order.IsReturned:
			return fmt.Errorf("can not give orders: order %d is already returned by recipient", id)
		case order.ExpireDate.Before(time.Now()):
			return fmt.Errorf("can not give orders: order %d is expired", id)
		case recipient == 0:
			recipient = order.RecipientID
		}
	}
	err = s.s.Give(ids)
	if err != nil {
		return err
	}

	return nil
}

// List checks validity of given recipient id and n and returns slice of all his orders (last n)
func (s service) List(ctx context.Context, recipient int64, n int, onlyNotGivenOrders bool) (_ []model.Order, err error) {
	defer func() {
		s.failedCounter.FailedCounterInc(err)
	}()
	if recipient <= 0 {
		return nil, errors.New("recipient id should be positive")
	}
	if n < 0 {
		return nil, errors.New("n should not be negative")
	}
	var (
		all []model.Order
	)
	if !onlyNotGivenOrders {
		all, err = s.s.ListAll(recipient)
	} else {
		all, err = s.s.ListNotGiven(recipient)
	}
	if err != nil {
		return nil, err
	}
	if n == 0 || len(all) <= n {
		return all, nil
	}

	return all[len(all)-n:], nil
}

// Return checks validity of given order id and recipient id and gets order back from recipient
func (s service) Return(ctx context.Context, id, recipient int64) (err error) {
	defer func() {
		s.failedCounter.FailedCounterInc(err)
	}()
	defer func() {
		s.givenOrdersGauge.SuccessGaugeDec(err)
	}()
	if id <= 0 {
		return errors.New("id should be positive")
	}
	if recipient <= 0 {
		return errors.New("recipient id should be positive")
	}
	order, ok := s.s.GetByID(id)
	switch {
	case !ok:
		return errors.New("order can not be returned: order not found")
	case order.RecipientID != recipient:
		return errors.New("order can not be returned: order belongs to different recipient")
	case order.IsReturned:
		return errors.New("order can not be returned: order is already returned")
	case !order.IsGiven:
		return errors.New("order can not be returned: order is not given yet")
	case order.GivenTime.AddDate(0, 0, 2).Before(time.Now()):
		return errors.New("order can not be returned: more than 2 days passed")
	}
	err = s.s.Return(id)
	if err != nil {
		return err
	}

	return nil
}

// ListReturn checks validity of given args and returns k returned orders on nth page
func (s service) ListReturn(ctx context.Context, pageNum, ordersPerPage int) (_ []model.Order, err error) {
	defer func() {
		s.failedCounter.FailedCounterInc(err)
	}()
	if pageNum < 0 {
		return nil, errors.New("pageNum should not be negative")
	}
	if ordersPerPage <= 0 {
		return nil, errors.New("ordersPerPage should be positive")
	}
	all, err := s.s.ListReturn()
	if err != nil {
		return nil, err
	}
	if pageNum == 0 {
		return all, nil
	}
	firstPos := (pageNum - 1) * ordersPerPage
	if len(all) == 0 || len(all) <= firstPos {
		return nil, errors.New("empty list")
	}
	newLen := ordersPerPage
	if len(all) < pageNum*ordersPerPage {
		newLen = len(all) % ordersPerPage
	}
	return all[firstPos : firstPos+newLen], nil
}
