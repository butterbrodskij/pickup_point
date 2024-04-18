//go:generate mockgen -source=./service.go -destination=./service_mocks_test.go -package=pickpoint
package pickpoint

import (
	"context"

	"gitlab.ozon.dev/mer_marat/homework/internal/model"
)

type storage interface {
	Add(context.Context, *model.PickPoint) (int64, error)
	GetByID(context.Context, int64) (*model.PickPoint, error)
	Update(context.Context, *model.PickPoint) error
	Delete(context.Context, int64) error
}

type cache interface {
	SetPickPoint(id int64, point model.PickPoint)
	GetPickPoint(id int64) (model.PickPoint, error)
	DeletePickPoint(id int64)
}

type service struct {
	repo  storage
	cache cache
}

// New returns type Service associated with storage
func NewService(stor storage, cache cache) service {
	return service{
		repo:  stor,
		cache: cache,
	}
}

// Create writes pick-up points information in storage
func (s service) Create(ctx context.Context, point *model.PickPoint) (*model.PickPoint, error) {
	id, err := s.repo.Add(ctx, point)
	if err != nil {
		return nil, err
	}
	point.ID = id
	return point, nil
}

// Read gets pick-up points information from storage by id
func (s service) Read(ctx context.Context, id int64) (*model.PickPoint, error) {
	if !isValidID(id) {
		return nil, model.ErrorInvalidInput
	}
	point, err := s.cache.GetPickPoint(id)
	if err == nil {
		return &point, nil
	}
	pPoint, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	s.cache.SetPickPoint(id, *pPoint)
	return pPoint, err
}

func (s service) Update(ctx context.Context, point *model.PickPoint) error {
	if !isValidPickPoint(point) {
		return model.ErrorInvalidInput
	}
	s.cache.DeletePickPoint(point.ID)
	return s.repo.Update(ctx, point)
}

func (s service) Delete(ctx context.Context, id int64) error {
	if !isValidID(id) {
		return model.ErrorInvalidInput
	}
	s.cache.DeletePickPoint(id)
	return s.repo.Delete(ctx, id)
}

func isValidPickPoint(point *model.PickPoint) bool {
	return point.ID > 0
}

func isValidID(id int64) bool {
	return id > 0
}
