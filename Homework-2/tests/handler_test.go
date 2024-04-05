package tests

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.ozon.dev/mer_marat/homework/internal/service/pickpoint"
	"gitlab.ozon.dev/mer_marat/homework/internal/storage/postgres"
	"gitlab.ozon.dev/mer_marat/homework/tests/fixture"
)

func TestCreate(t *testing.T) {
	var (
		ctx = context.Background()
	)
	t.Run("smoke test", func(t *testing.T) {
		db.SetUp(t, "pickpoints")
		defer db.TearDown()
		repo := postgres.NewRepo(db.DB)
		serv := pickpoint.NewService(repo)

		res, err := serv.Create(ctx, fixture.PickPoint().Valid1().P())

		require.NoError(t, err)
		assert.Equal(t, res.ID, int64(1))
	})
}
