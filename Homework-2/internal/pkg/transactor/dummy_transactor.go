// needs for file storage where transactions not supported
package transactor

import (
	"context"

	"github.com/jackc/pgx/v4"
)

type DummyTransactor struct {
}

func NewDummyTransactor() *DummyTransactor {
	return &DummyTransactor{}
}

func (t *DummyTransactor) RunSerializable(ctx context.Context, role pgx.TxAccessMode, f func(ctxTX context.Context) error) error {
	return f(ctx)
}
