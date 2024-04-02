package postgres

import (
	"context"
	"errors"

	"gitlab.ozon.dev/mer_marat/homework/internal/model"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

type Database interface {
	Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error)
	ExecQueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row
}

type PickpointRepo struct {
	db Database
}

func NewRepo(db Database) *PickpointRepo {
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

func (r *PickpointRepo) Update(ctx context.Context, point *model.PickPoint) error {
	query := "UPDATE pickpoints SET name=$1, address=$2, contacts=$3 WHERE id=$4"
	tag, err := r.db.Exec(ctx, query, point.Name, point.Address, point.Contact, point.ID)
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
	tag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return model.ErrorObjectNotFound
	}
	return nil
}
