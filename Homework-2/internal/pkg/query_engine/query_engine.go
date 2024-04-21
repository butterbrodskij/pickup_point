package queryengine

import (
	"context"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

type QueryEngine interface {
	Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error)
	ExecQueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row
}

type QueryEngineTx struct {
	tx pgx.Tx
}

func NewQueryEngineTx(tx pgx.Tx) *QueryEngineTx {
	return &QueryEngineTx{
		tx: tx,
	}
}

func (q *QueryEngineTx) Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return pgxscan.Get(ctx, q.tx, dest, query, args...)
}

func (q *QueryEngineTx) Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return pgxscan.Select(ctx, q.tx, dest, query, args...)
}

func (q *QueryEngineTx) Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error) {
	return q.tx.Exec(ctx, query, args...)
}

func (q *QueryEngineTx) ExecQueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row {
	return q.tx.QueryRow(ctx, query, args...)
}
