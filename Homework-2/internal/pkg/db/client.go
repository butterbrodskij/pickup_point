package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
	"gitlab.ozon.dev/mer_marat/homework/cmd/config"
)

func NewDB(ctx context.Context, cfg config.Config) (*Database, error) {
	pool, err := pgxpool.Connect(ctx, generateDsn(cfg))
	if err != nil {
		return nil, err
	}
	return newDatabase(pool), nil
}

func generateDsn(cfg config.Config) string {
	pattern := "host=%s port=%d user=%s password=%s dbname=%s sslmode=disable"
	return fmt.Sprintf(pattern, cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, cfg.Database.Dbname)
}
