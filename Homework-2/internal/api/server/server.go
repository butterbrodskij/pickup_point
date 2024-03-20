package server

import (
	"context"

	"gitlab.ozon.dev/mer_marat/homework/internal/model"
	"gitlab.ozon.dev/mer_marat/homework/internal/pkg/repository/postgres"

	"github.com/jackc/pgconn"
)

type Server struct {
	repo *postgres.PickpointRepo
}

func NewServer(repo *postgres.PickpointRepo) Server {
	return Server{repo: repo}
}

func (s Server) Create(ctx context.Context, point *model.PickPoint) (int64, error) {
	return s.repo.Add(ctx, point)
}

func (s Server) Update(ctx context.Context, point *model.PickPoint) (pgconn.CommandTag, error) {
	return s.repo.Update(ctx, point)
}

func (s Server) Read(ctx context.Context, id int64) (*model.PickPoint, error) {
	return s.repo.GetByID(ctx, id)
}

func (s Server) Delete(ctx context.Context, id int64) (pgconn.CommandTag, error) {
	return s.repo.Delete(ctx, id)
}
