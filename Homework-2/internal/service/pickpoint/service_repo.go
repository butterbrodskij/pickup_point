package pickpoint

import (
	"context"
	"strconv"

	"github.com/jackc/pgconn"
	"gitlab.ozon.dev/mer_marat/homework/internal/model"
)

type repoInterface interface {
	Add(context.Context, *model.PickPoint) (int64, error)
	GetByID(context.Context, int64) (*model.PickPoint, error)
	Update(context.Context, *model.PickPoint) (pgconn.CommandTag, error)
	Delete(context.Context, int64) (pgconn.CommandTag, error)
}

type ServiceRepo struct {
	repo repoInterface
}

func NewServiceRepo(repo repoInterface) ServiceRepo {
	return ServiceRepo{repo: repo}
}

func (s ServiceRepo) Create(ctx context.Context, point *model.PickPoint) (*model.PickPoint, error) {
	id, err := s.repo.Add(ctx, point)
	if err != nil {
		return nil, err
	}
	point.ID = id
	return point, nil
}

func (s ServiceRepo) Update(ctx context.Context, point *model.PickPoint) (pgconn.CommandTag, error) {
	if !validPickPoint(point) {
		return nil, model.ErrorInvalidInput
	}
	return s.repo.Update(ctx, point)
}

func (s ServiceRepo) Read(ctx context.Context, ids string) (*model.PickPoint, error) {
	id, err := strconv.ParseInt(ids, 10, 64)
	if err != nil {
		return nil, err
	}
	if !validID(id) {
		return nil, model.ErrorInvalidInput
	}
	return s.repo.GetByID(ctx, id)
}

func (s ServiceRepo) Delete(ctx context.Context, ids string) (pgconn.CommandTag, error) {
	id, err := strconv.ParseInt(ids, 10, 64)
	if err != nil {
		return nil, err
	}
	if !validID(id) {
		return nil, model.ErrorInvalidInput
	}
	return s.repo.Delete(ctx, id)
}

func validPickPoint(point *model.PickPoint) bool {
	return point.ID > 0
}

func validID(id int64) bool {
	return id > 0
}
