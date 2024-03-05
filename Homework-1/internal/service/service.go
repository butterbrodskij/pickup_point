package service

import (
	"errors"
	"homework1/pup/internal/model"
	"strconv"
	"time"
)

type storage interface {
	Get(model.Order) error
	Remove(int) error
	Give([]int) error
	List(int, bool) ([]model.Order, error)
}

type Service struct {
	s storage
}

func Input2Order(input model.OrderInput) (model.Order, error) {
	if input.ID <= 0 {
		return model.Order{}, errors.New("id should be positive")
	}

	if input.Recipient <= 0 {
		return model.Order{}, errors.New("recipient id should be positive")
	}

	t, err := time.Parse("2.1.2006", input.ExpireDate)
	if err != nil {
		return model.Order{}, errors.New("wrong date format")
	}

	if t.Before(time.Now()) {
		return model.Order{}, errors.New("trying to get expired order")
	}

	return model.Order{
		ID:         input.ID,
		Recipient:  input.Recipient,
		ExpireDate: t,
	}, nil
}

func New(stor storage) Service {
	return Service{s: stor}
}

func (s Service) Get(input model.OrderInput) error {
	order, err := Input2Order(input)
	if err != nil {
		return err
	}
	return s.s.Get(order)
}

func (s Service) Remove(id int) error {
	if id <= 0 {
		return errors.New("id should be positive")
	}
	return s.s.Remove(id)
}

func (s Service) Give(idString []string) error {
	ids := make([]int, len(idString))
	for i, str := range idString {
		id, err := strconv.Atoi(str)
		if err != nil {
			return err
		}
		ids[i] = id
	}
	return s.s.Give(ids)
}

func (s Service) List(recipient, n int, flag bool) ([]model.Order, error) {
	if recipient <= 0 {
		return []model.Order{}, errors.New("recipient id should be positive")
	}
	if n < 0 {
		return []model.Order{}, errors.New("n should not be negative")
	}
	all, err := s.s.List(recipient, flag)
	if err != nil || n == 0 || len(all) <= n {
		return all, err
	}

	return all[:n], err
}
