package transactor

import (
	"context"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	queryengine "gitlab.ozon.dev/mer_marat/homework/internal/pkg/query_engine"
)

type keyType string

const key keyType = "transaction"

type database interface {
	Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error)
	ExecQueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row
	BeginTx(ctx context.Context, opt pgx.TxOptions) (pgx.Tx, error)
}

type Transactor struct {
	pool database
}

func NewTransactor(pool database) *Transactor {
	return &Transactor{
		pool: pool,
	}
}

func (t *Transactor) RunSerializable(ctx context.Context, f func(ctxTX context.Context) error) error {
	tx, err := t.pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   pgx.Serializable,
		AccessMode: pgx.ReadWrite,
	})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if err := f(context.WithValue(ctx, key, queryengine.NewQueryEngineTx(tx))); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (t *Transactor) GetQueryEngine(ctx context.Context) queryengine.QueryEngine {
	tx, ok := ctx.Value(key).(queryengine.QueryEngine)
	if ok {
		return tx
	}
	return t.pool
}
