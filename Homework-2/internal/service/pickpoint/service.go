package pickpoint

import (
	"context"

	"gitlab.ozon.dev/mer_marat/homework/internal/model"
)

type storageInterface interface {
	Add(context.Context, *model.PickPoint) (int64, error)
	GetByID(context.Context, int64) (*model.PickPoint, error)
}

type ServiceInterface interface {
	Create(context.Context, *model.PickPoint) (*model.PickPoint, error)
	Read(context.Context, int64) (*model.PickPoint, error)
}

type Service struct {
	repo storageInterface
}

// New returns type Service associated with storage
func New(stor storageInterface) Service {
	return Service{repo: stor}
}

// Create writes pick-up points information in storage
func (s Service) Create(ctx context.Context, point *model.PickPoint) (*model.PickPoint, error) {
	id, err := s.repo.Add(ctx, point)
	if err != nil {
		return nil, err
	}
	point.ID = id
	return point, nil
}

// Read gets pick-up points information from storage by id
func (s Service) Read(ctx context.Context, id int64) (*model.PickPoint, error) {
	if !validID(id) {
		return nil, model.ErrorInvalidInput
	}
	return s.repo.GetByID(ctx, id)
}
