//go:build !integration

package postgres

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.ozon.dev/mer_marat/homework/internal/model"
)

func TestAdd(t *testing.T) {
	t.Parallel()
}

func TestGetByID(t *testing.T) {
	t.Parallel()
	var (
		ctx = context.Background()
		id  = int64(1)
	)
	t.Run("smoke test", func(t *testing.T) {
		t.Parallel()
		s := setUp(t)
		defer s.tearDown()
		s.mockDB.EXPECT().Get(gomock.Any(), gomock.Any(), "SELECT id, name, address, contacts FROM pickpoints WHERE id=$1", gomock.Any()).Return(nil)

		user, err := s.repo.GetByID(ctx, id)

		require.NoError(t, err)
		assert.Equal(t, int64(0), user.ID)
	})
	t.Run("fail", func(t *testing.T) {
		t.Parallel()
		t.Run("not found", func(t *testing.T) {
			t.Parallel()
			s := setUp(t)
			defer s.tearDown()
			s.mockDB.EXPECT().Get(gomock.Any(), gomock.Any(), "SELECT id, name, address, contacts FROM pickpoints WHERE id=$1", gomock.Any()).Return(pgx.ErrNoRows)

			user, err := s.repo.GetByID(ctx, id)

			require.EqualError(t, err, "object not found")
			require.True(t, errors.Is(err, model.ErrorObjectNotFound))
			assert.Nil(t, user)
		})
		t.Run("internal error", func(t *testing.T) {
			t.Parallel()
			s := setUp(t)
			defer s.tearDown()
			s.mockDB.EXPECT().Get(gomock.Any(), gomock.Any(), "SELECT id, name, address, contacts FROM pickpoints WHERE id=$1", gomock.Any()).Return(assert.AnError)

			user, err := s.repo.GetByID(ctx, id)

			require.EqualError(t, err, "assert.AnError general error for testing")
			assert.Nil(t, user)
		})
	})
}

func TestUpdate(t *testing.T) {
	t.Parallel()
}

func TestDelete(t *testing.T) {
	t.Parallel()
}
