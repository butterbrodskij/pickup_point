package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
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
	)
	t.Run("create changes id test", func(t *testing.T) {
		t.Parallel()
		s := setUp(t)
		defer s.tearDown()
		body := `{"name":"Chertanovo", "address":"Chertanovskaya street, 8", "contacts":"+7(999)888-77-66"}`
		req, _ := http.NewRequestWithContext(ctx, "POST", "/pickpoint", strings.NewReader(body))
		w := httptest.NewRecorder()
		s.mockServ.EXPECT().Create(gomock.Any(), gomock.Any()).Return(fixture.PickPoint().Valid1().P(), nil)

		s.handl.Create(w, req)

		require.Equal(t, w.Code, http.StatusOK)
		assert.Equal(t, w.Body.String(), `{"id":100,"name":"Chertanovo","address":"Chertanovskaya street, 8","contacts":"+7(999)888-77-66"}`)
	})
	t.Run("fail", func(t *testing.T) {
		t.Parallel()
		t.Run("bad request", func(t *testing.T) {
			t.Parallel()
			s := setUp(t)
			defer s.tearDown()
			body := `example of bad request`
			req, _ := http.NewRequestWithContext(ctx, "POST", "/pickpoint", strings.NewReader(body))
			w := httptest.NewRecorder()

			s.handl.Create(w, req)

			require.Equal(t, w.Code, http.StatusBadRequest)
			assert.Equal(t, w.Body.String(), "")
		})
		t.Run("internal error", func(t *testing.T) {
			t.Parallel()
			s := setUp(t)
			defer s.tearDown()
			body := `{"name":"Chertanovo", "address":"Chertanovskaya street, 13", "contacts":"+7(999)888-77-66"}`
			req, _ := http.NewRequestWithContext(ctx, "POST", "/pickpoint", strings.NewReader(body))
			w := httptest.NewRecorder()
			s.mockServ.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil, assert.AnError)

			s.handl.Create(w, req)

			require.Equal(t, w.Code, http.StatusInternalServerError)
			assert.Equal(t, w.Body.String(), "")
		})
	})
}
