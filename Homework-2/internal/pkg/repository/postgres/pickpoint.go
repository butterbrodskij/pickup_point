package postgres

import (
	"context"
	"database/sql"
	"errors"
	"homework2/pup/internal/model"
	"homework2/pup/internal/pkg/db"
)

type PickpointRepo struct {
	db *db.Database
}

func NewRepo(db *db.Database) *PickpointRepo {
	return &PickpointRepo{db: db}
}

func (r *PickpointRepo) Add(ctx context.Context) (int64, error) {
	return 0, nil
}

func (r *PickpointRepo) GetByID(ctx context.Context, id int64) (*model.PickPoint, error) {
	var point model.PickPoint
	query := "SELECT id, name, address, contacts FROM pickpoints WHERE id=$1"
	err := r.db.Get(ctx, &point, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrorObjectNotFound
		}
		return nil, err
	}
	return &point, nil
}

func (r *PickpointRepo) Delete(ctx context.Context) error {
	return nil
}
