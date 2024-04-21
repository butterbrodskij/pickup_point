//go:generate mockgen -source=./postgres_point.go -destination=./postgres_point_mocks_test.go -package=postgres
package postgres

import (
	"context"
	"errors"

	"gitlab.ozon.dev/mer_marat/homework/internal/model"
	queryengine "gitlab.ozon.dev/mer_marat/homework/internal/pkg/query_engine"

	"github.com/jackc/pgx/v4"
)

type transactor interface {
	GetQueryEngine(ctx context.Context) queryengine.QueryEngine
	RunSerializable(ctx context.Context, f func(ctxTX context.Context) error) error
}

type PickpointRepo struct {
	db transactor
}

func NewRepo(db transactor) *PickpointRepo {
	return &PickpointRepo{db: db}
}

func (r *PickpointRepo) Add(ctx context.Context, point *model.PickPoint) (int64, error) {
	var id int64
	query := "INSERT INTO pickpoints(name, address, contacts) VALUES ($1, $2, $3) RETURNING id;"
	err := r.db.GetQueryEngine(ctx).ExecQueryRow(ctx, query, point.Name, point.Address, point.Contact).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *PickpointRepo) GetByID(ctx context.Context, id int64) (*model.PickPoint, error) {
	var point model.PickPoint
	query := "SELECT id, name, address, contacts FROM pickpoints WHERE id=$1"
	err := r.db.GetQueryEngine(ctx).Get(ctx, &point, query, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrorObjectNotFound
		}
		return nil, err
	}
	return &point, nil
}

func (r *PickpointRepo) Update(ctx context.Context, point *model.PickPoint) error {
	query := "UPDATE pickpoints SET name=$1, address=$2, contacts=$3 WHERE id=$4"
	tag, err := r.db.GetQueryEngine(ctx).Exec(ctx, query, point.Name, point.Address, point.Contact, point.ID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return model.ErrorObjectNotFound
	}
	return nil
}

func (r *PickpointRepo) Delete(ctx context.Context, id int64) error {
	query := "DELETE FROM pickpoints WHERE id=$1"
	tag, err := r.db.GetQueryEngine(ctx).Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return model.ErrorObjectNotFound
	}
	return nil
}
