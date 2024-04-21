package pickpoint

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.ozon.dev/mer_marat/homework/internal/model"
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

func TestRead(t *testing.T) {
	t.Parallel()
	var (
		ctx       = context.Background()
		id        = int64(1)
		invalidID = int64(-1)
	)
	t.Run("smoke test", func(t *testing.T) {
		t.Parallel()
		s := setUp(t)
		defer s.tearDown()
		point := fixture.PickPoint().Valid1()
		s.mockCache.EXPECT().Get(gomock.Any(), fmt.Sprint(id), gomock.Any()).Return(model.ErrorCacheMissed)
		s.mockRepo.EXPECT().GetByID(gomock.Any(), id).Return(point.P(), nil)
		s.mockCache.EXPECT().Set(gomock.Any(), fmt.Sprint(id), gomock.Any()).Return(nil)

		result, err := s.serv.Read(ctx, id)

		require.NoError(t, err)
		assert.Equal(t, point.P(), result)
	})
	t.Run("cache successful", func(t *testing.T) {
		t.Parallel()
		s := setUp(t)
		defer s.tearDown()
		s.mockCache.EXPECT().Get(gomock.Any(), fmt.Sprint(id), gomock.Any()).Return(nil)

		_, err := s.serv.Read(ctx, id)

		require.NoError(t, err)
	})
	t.Run("fail", func(t *testing.T) {
		t.Parallel()
		t.Run("error: invalid input", func(t *testing.T) {
			t.Parallel()
			s := setUp(t)
			defer s.tearDown()

			result, err := s.serv.Read(ctx, invalidID)

			require.EqualError(t, err, "invalid input")
			assert.Nil(t, result)
		})
		t.Run("error", func(t *testing.T) {
			t.Parallel()
			s := setUp(t)
			defer s.tearDown()
			s.mockCache.EXPECT().Get(gomock.Any(), fmt.Sprint(id), gomock.Any()).Return(model.ErrorCacheMissed)
			s.mockRepo.EXPECT().GetByID(gomock.Any(), id).Return(nil, assert.AnError)

			result, err := s.serv.Read(ctx, id)

			require.EqualError(t, err, "assert.AnError general error for testing")
			assert.Nil(t, result)
		})
	})
}

func TestUpdate(t *testing.T) {
	t.Parallel()
	var (
		ctx = context.Background()
	)
	t.Run("smoke test", func(t *testing.T) {
		t.Parallel()
		s := setUp(t)
		defer s.tearDown()
		point := fixture.PickPoint().Valid1()
		s.mockRepo.EXPECT().Update(gomock.Any(), point.P()).Return(nil)
		s.mockCache.EXPECT().Delete(gomock.Any(), fmt.Sprint(point.V().ID)).Return(nil)

		err := s.serv.Update(ctx, point.P())

		require.NoError(t, err)
	})
	t.Run("fail", func(t *testing.T) {
		t.Parallel()
		t.Run("error: invalid input", func(t *testing.T) {
			t.Parallel()
			s := setUp(t)
			defer s.tearDown()
			invalidPoint := fixture.PickPoint().InValid1()

			err := s.serv.Update(ctx, invalidPoint.P())

			require.EqualError(t, err, "invalid input")
		})
		t.Run("error", func(t *testing.T) {
			t.Parallel()
			s := setUp(t)
			defer s.tearDown()
			point := fixture.PickPoint().Valid1()
			s.mockCache.EXPECT().Delete(gomock.Any(), fmt.Sprint(point.V().ID)).Return(nil)
			s.mockRepo.EXPECT().Update(gomock.Any(), point.P()).Return(assert.AnError)

			err := s.serv.Update(ctx, point.P())

			require.EqualError(t, err, "assert.AnError general error for testing")
		})
	})
}

func TestDelete(t *testing.T) {
	t.Parallel()
	var (
		ctx       = context.Background()
		id        = int64(1)
		invalidID = int64(-1)
	)
	t.Run("smoke test", func(t *testing.T) {
		t.Parallel()
		s := setUp(t)
		defer s.tearDown()
		s.mockCache.EXPECT().Delete(gomock.Any(), fmt.Sprint(id)).Return(nil)
		s.mockRepo.EXPECT().Delete(gomock.Any(), id).Return(nil)

		err := s.serv.Delete(ctx, id)

		require.NoError(t, err)
	})
	t.Run("fail", func(t *testing.T) {
		t.Parallel()
		t.Run("error: invalid input", func(t *testing.T) {
			t.Parallel()
			s := setUp(t)
			defer s.tearDown()

			err := s.serv.Delete(ctx, invalidID)

			require.EqualError(t, err, "invalid input")
		})
		t.Run("error", func(t *testing.T) {
			t.Parallel()
			s := setUp(t)
			defer s.tearDown()
			s.mockCache.EXPECT().Delete(gomock.Any(), fmt.Sprint(id)).Return(nil)
			s.mockRepo.EXPECT().Delete(gomock.Any(), id).Return(assert.AnError)

			err := s.serv.Delete(ctx, id)

			require.EqualError(t, err, "assert.AnError general error for testing")
		})
	})
}
