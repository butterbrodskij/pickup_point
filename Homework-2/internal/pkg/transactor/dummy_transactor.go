// needs for file storage where transactions not supported
package transactor

import (
	"context"
)

type DummyTransactor struct {
}

func NewDummyTransactor() *DummyTransactor {
	return &DummyTransactor{}
}

func (t *DummyTransactor) RunSerializable(ctx context.Context, f func(ctxTX context.Context) error) error {
	return f(ctx)
}
