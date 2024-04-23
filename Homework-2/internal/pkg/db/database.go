package db

import (
	"context"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgtype/pgxtype"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type keyType string

const key keyType = "transaction"

type Database struct {
	cluster *pgxpool.Pool
}

func newDatabase(cluster *pgxpool.Pool) *Database {
	return &Database{cluster: cluster}
}

func (db Database) Close() {
	db.cluster.Close()
}

func (db Database) Pool() *pgxpool.Pool {
	return db.cluster
}

func (db Database) BeginTx(ctx context.Context, opt pgx.TxOptions) (pgx.Tx, error) {
	return db.cluster.BeginTx(ctx, opt)
}

func (db Database) RunSerializable(ctx context.Context, f func(ctxTX context.Context) error) error {
	tx, ok := ctx.Value(key).(pgx.Tx)
	if !ok {
		var err error
		tx, err = db.BeginTx(ctx, pgx.TxOptions{
			IsoLevel:   pgx.Serializable,
			AccessMode: pgx.ReadWrite,
		})
		if err != nil {
			return err
		}
	}
	defer tx.Rollback(ctx)
	if err := f(context.WithValue(ctx, key, tx)); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (db Database) GetQuerier(ctx context.Context) pgxtype.Querier {
	tx, ok := ctx.Value(key).(pgxtype.Querier)
	if ok {
		return tx
	}
	return db.cluster
}

func (db Database) Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return pgxscan.Get(ctx, db.GetQuerier(ctx), dest, query, args...)
}

func (db Database) Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return pgxscan.Select(ctx, db.GetQuerier(ctx), dest, query, args...)
}

func (db Database) Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error) {
	return db.GetQuerier(ctx).Exec(ctx, query, args...)
}

func (db Database) ExecQueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row {
	return db.GetQuerier(ctx).QueryRow(ctx, query, args...)
}
