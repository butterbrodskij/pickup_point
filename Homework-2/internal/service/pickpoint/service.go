//go:generate mockgen -source=./service.go -destination=./service_mocks_test.go -package=pickpoint
package pickpoint

import (
	"context"
	"fmt"
	"log"

	"gitlab.ozon.dev/mer_marat/homework/internal/model"
)

type storage interface {
	Add(context.Context, *model.PickPoint) (int64, error)
	GetByID(context.Context, int64) (*model.PickPoint, error)
	Update(context.Context, *model.PickPoint) error
	Delete(context.Context, int64) error
}

type cache interface {
	Set(ctx context.Context, key string, value interface{}) error
	Get(ctx context.Context, key string, value interface{}) error
	Delete(ctx context.Context, keys ...string) error
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
	point := new(model.PickPoint)
	err := s.cache.Get(ctx, fmt.Sprint(id), point)
	if err == nil {
		return point, nil
	}
	pPoint, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	err = s.cache.Set(ctx, fmt.Sprint(id), *pPoint)
	if err != nil {
		log.Printf("cache set failed: %s", err)
	}
	return pPoint, nil
}

func (s service) Update(ctx context.Context, point *model.PickPoint) error {
	if !isValidPickPoint(point) {
		return model.ErrorInvalidInput
	}
	err := s.cache.Delete(ctx, fmt.Sprint(point.ID))
	if err != nil {
		return err
	}
	return s.repo.Update(ctx, point)
}

func (s service) Delete(ctx context.Context, id int64) error {
	if !isValidID(id) {
		return model.ErrorInvalidInput
	}
	err := s.cache.Delete(ctx, fmt.Sprint(id))
	if err != nil {
		return err
	}
	return s.repo.Delete(ctx, id)
}

func isValidPickPoint(point *model.PickPoint) bool {
	return point.ID > 0
}

func isValidID(id int64) bool {
	return id > 0
}
