package postgres

import (
	"context"
	"errors"
	"homework2/pup/internal/model"
	"homework2/pup/internal/pkg/db"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

type PickpointRepo struct {
	db *db.Database
}

func NewRepo(db *db.Database) *PickpointRepo {
	return &PickpointRepo{db: db}
}

func (r *PickpointRepo) Add(ctx context.Context, point *model.PickPoint) (int64, error) {
	var id int64
	query := "INSERT INTO pickpoints(name, address, contacts) VALUES ($1, $2, $3) RETURNING id;"
	err := r.db.ExecQueryRow(ctx, query, point.Name, point.Address, point.Contact).Scan(&id)
	return id, err
}

func (r *PickpointRepo) GetByID(ctx context.Context, id int64) (*model.PickPoint, error) {
	var point model.PickPoint
	query := "SELECT id, name, address, contacts FROM pickpoints WHERE id=$1"
	err := r.db.Get(ctx, &point, query, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrorObjectNotFound
		}
		return nil, err
	}
	return &point, nil
}

func (r *PickpointRepo) Update(ctx context.Context, point *model.PickPoint) (pgconn.CommandTag, error) {
	query := "UPDATE pickpoints SET name=$1, address=$2, contacts=$3 WHERE id=$4"
	return r.db.Exec(ctx, query, point.Name, point.Address, point.Contact, point.ID)
}

func (r *PickpointRepo) Delete(ctx context.Context, id int64) (pgconn.CommandTag, error) {
	query := "DELETE FROM pickpoints WHERE id=$1"
	return r.db.Exec(ctx, query, id)
}
