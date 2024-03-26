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
