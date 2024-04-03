package pickpoint

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.ozon.dev/mer_marat/homework/tests/fixture"
)

func TestCreate(t *testing.T) {
	t.Parallel()
	var (
		ctx = context.Background()
		id  = int64(1)
	)
	t.Run("create changes id test", func(t *testing.T) {
		t.Parallel()
		s := setUp(t)
		defer s.tearDown()
		point := fixture.PickPoint().Valid1()
		s.mockRepo.EXPECT().Add(gomock.Any(), point.P()).Return(id, nil)

		result, err := s.serv.Create(ctx, point.P())

		require.NoError(t, err)
		assert.Equal(t, id, result.ID)
	})
	t.Run("error", func(t *testing.T) {
		t.Parallel()
		s := setUp(t)
		defer s.tearDown()
		point := fixture.PickPoint().Valid1()
		s.mockRepo.EXPECT().Add(gomock.Any(), point.P()).Return(id, assert.AnError)

		result, err := s.serv.Create(ctx, point.P())

		require.EqualError(t, err, "assert.AnError general error for testing")
		assert.Nil(t, result)
	})
}
