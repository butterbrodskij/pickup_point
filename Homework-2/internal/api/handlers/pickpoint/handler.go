//go:generate mockgen -source=./handler.go -destination=./handler_mocks_test.go -package=handler
package handler

import (
	"context"

	"gitlab.ozon.dev/mer_marat/homework/internal/model"
)

type service interface {
	Create(context.Context, *model.PickPoint) (*model.PickPoint, error)
	Read(context.Context, int64) (*model.PickPoint, error)
	Update(context.Context, *model.PickPoint) error
	Delete(context.Context, int64) error
}

type cache interface {
	Set(ctx context.Context, key string, value interface{}) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key ...string) error
}

type handler struct {
	service
	cache
}

func NewHandler(s service, cache cache) *handler {
	return &handler{
		service: s,
		cache:   cache,
	}
}
