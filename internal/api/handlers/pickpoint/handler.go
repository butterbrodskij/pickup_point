//go:generate mockgen -source=./handler.go -destination=./handler_mocks_test.go -package=handler
package handler

import (
	"context"

	"gitlab.ozon.dev/mer_marat/homework/internal/model"
)

type service interface {
	Read(ctx context.Context, id int64) (res *model.PickPoint, err error)
	Create(ctx context.Context, point *model.PickPoint) (res *model.PickPoint, err error)
	Update(ctx context.Context, point *model.PickPoint) (err error)
	Delete(ctx context.Context, id int64) (err error)
}

type handler struct {
	service
}

func NewHandler(s service) *handler {
	return &handler{
		service: s,
	}
}
