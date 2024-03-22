package pickpoint

import (
	"context"

	"github.com/jackc/pgconn"
	"gitlab.ozon.dev/mer_marat/homework/internal/model"
)

type repoInterface interface {
	storageInterface
	Update(context.Context, *model.PickPoint) (pgconn.CommandTag, error)
	Delete(context.Context, int64) (pgconn.CommandTag, error)
}

type ServiceRepoInteface interface {
	ServiceInterface
	Update(context.Context, *model.PickPoint) (pgconn.CommandTag, error)
	Delete(context.Context, int64) (pgconn.CommandTag, error)
}

type ServiceRepo struct {
	ServiceInterface
	repo repoInterface
}

func NewServiceRepo(repo repoInterface) ServiceRepo {
	return ServiceRepo{ServiceInterface: Service{repo: repo}, repo: repo}
}

func (s ServiceRepo) Update(ctx context.Context, point *model.PickPoint) (pgconn.CommandTag, error) {
	if !validPickPoint(point) {
		return nil, model.ErrorInvalidInput
	}
	return s.repo.Update(ctx, point)
}

func (s ServiceRepo) Delete(ctx context.Context, id int64) (pgconn.CommandTag, error) {
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
