//go:generate mockgen -source=./postgres_point.go -destination=./postgres_point_mocks_test.go -package=postgres
package postgres

import (
	"context"
	"errors"

	"gitlab.ozon.dev/mer_marat/homework/internal/model"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

type database interface {
	Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error)
	ExecQueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row
	BeginTx(ctx context.Context, opt pgx.TxOptions) (pgx.Tx, error)
}

type PickpointRepo struct {
	db database
}

func NewRepo(db database) *PickpointRepo {
	return &PickpointRepo{db: db}
}

func (r *PickpointRepo) Add(ctx context.Context, point *model.PickPoint) (int64, error) {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)
	var id int64
	query := "INSERT INTO pickpoints(name, address, contacts) VALUES ($1, $2, $3) RETURNING id;"
	err = r.db.ExecQueryRow(ctx, query, point.Name, point.Address, point.Contact).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, tx.Commit(ctx)
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
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	query := "UPDATE pickpoints SET name=$1, address=$2, contacts=$3 WHERE id=$4"
	tag, err := r.db.Exec(ctx, query, point.Name, point.Address, point.Contact, point.ID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return model.ErrorObjectNotFound
	}
	return tx.Commit(ctx)
}

func (r *PickpointRepo) Delete(ctx context.Context, id int64) error {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	query := "DELETE FROM pickpoints WHERE id=$1"
	tag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return model.ErrorObjectNotFound
	}
	return tx.Commit(ctx)
}
