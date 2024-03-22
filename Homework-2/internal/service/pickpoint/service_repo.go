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

func (s ServiceRepo) Create(ctx context.Context, point model.PickPointAdd) (*model.PickPoint, error) {
	pointRepo := &model.PickPoint{
		Name:    point.Name,
		Address: point.Address,
		Contact: point.Contact,
	}
	id, err := s.repo.Add(ctx, pointRepo)
	if err != nil {
		return nil, err
	}
	pointRepo.ID = id
	return pointRepo, nil
}

func (s ServiceRepo) Update(ctx context.Context, point *model.PickPoint) (pgconn.CommandTag, error) {
	return s.repo.Update(ctx, point)
}

func (s ServiceRepo) Read(ctx context.Context, ids string) (*model.PickPoint, error) {
	id, err := strconv.ParseInt(ids, 10, 64)
	if err != nil {
		return nil, err
	}
	return s.repo.GetByID(ctx, id)
}

func (s ServiceRepo) Delete(ctx context.Context, ids string) (pgconn.CommandTag, error) {
	id, err := strconv.ParseInt(ids, 10, 64)
	if err != nil {
		return nil, err
	}
	return s.repo.Delete(ctx, id)
}
