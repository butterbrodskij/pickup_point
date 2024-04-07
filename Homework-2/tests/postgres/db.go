package postgres_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"gitlab.ozon.dev/mer_marat/homework/internal/config"
	"gitlab.ozon.dev/mer_marat/homework/internal/pkg/db"
)

type database interface {
	Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error)
	ExecQueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row
}

type TDB struct {
	DB database
}

func NewTDB(cfg config.Config) *TDB {
	db, err := db.NewDB(context.Background(), cfg)
	if err != nil {
		panic(err)
	}
	return &TDB{DB: db}
}

func (d *TDB) SetUp(t *testing.T, tableName ...string) {
	t.Helper()
	d.truncateTable(context.Background(), tableName...)
}

func (d *TDB) TearDown(t *testing.T, tableName ...string) {
	t.Helper()
	d.truncateTable(context.Background(), tableName...)
}

func (d *TDB) truncateTable(ctx context.Context, tableName ...string) {
	q := fmt.Sprintf("TRUNCATE table %s RESTART IDENTITY;", strings.Join(tableName, ","))
	if _, err := d.DB.Exec(ctx, q); err != nil {
		panic(err)
	}
}
