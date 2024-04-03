//go:generate mockgen -source=./service.go -destination=./mocks/service.go -package=mock_service
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

type service struct {
	repo storage
}

// New returns type Service associated with storage
func NewService(stor storage) service {
	return service{repo: stor}
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
	return s.repo.GetByID(ctx, id)
}

func (s service) Update(ctx context.Context, point *model.PickPoint) error {
	if !isValidPickPoint(point) {
		return model.ErrorInvalidInput
	}
	return s.repo.Update(ctx, point)
}

func (s service) Delete(ctx context.Context, id int64) error {
	if !isValidID(id) {
		return model.ErrorInvalidInput
	}
	return s.repo.Delete(ctx, id)
}

func isValidPickPoint(point *model.PickPoint) bool {
	return point.ID > 0
}

func isValidID(id int64) bool {
	return id > 0
}
