package pickpoint

import (
	"fmt"

	"gitlab.ozon.dev/mer_marat/homework/internal/model"
)

type storageInterface interface {
	Write(model.PickPoint) error
	Get(int64) (model.PickPoint, bool)
}

type Service struct {
	s storageInterface
}

// New returns type Service associated with storage
func New(stor storageInterface) Service {
	return Service{s: stor}
}

// Write writes pick-up points information in storage
func (s Service) Write(point model.PickPoint) error {
	return s.s.Write(point)
}

// Get gets pick-up points information from storage by id
func (s Service) Get(id int64) (model.PickPoint, error) {
	point, ok := s.s.Get(id)
	if !ok {
		return model.PickPoint{}, fmt.Errorf("can not get point %d: not found", id)
	} else {
		return point, nil
	}
}
