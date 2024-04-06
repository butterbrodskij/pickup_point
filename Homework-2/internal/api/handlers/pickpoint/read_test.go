//go:build !integration

package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.ozon.dev/mer_marat/homework/internal/model"
	"gitlab.ozon.dev/mer_marat/homework/tests/fixture"
)

func TestRead(t *testing.T) {
	t.Parallel()
	var (
		ctx       = context.Background()
		id        = int64(100)
		invalidID = int64(0)
	)
	t.Run("smoke test", func(t *testing.T) {
		t.Parallel()
		s := setUp(t)
		defer s.tearDown()
		req, _ := http.NewRequestWithContext(ctx, "GET", "/pickpoint/100", strings.NewReader(""))
		w := httptest.NewRecorder()
		m := mux.NewRouter()
		s.mockServ.EXPECT().Read(gomock.Any(), id).Return(fixture.PickPoint().Valid1().P(), nil)
		m.HandleFunc("/pickpoint/{point:[0-9]+}", s.handl.Read)
		m.ServeHTTP(w, req)

		s.handl.Read(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, w.Body.String(), `{"id":100,"name":"Chertanovo","address":"Chertanovskaya street, 8","contacts":"+7(999)888-77-66"}`)
	})
	t.Run("fail", func(t *testing.T) {
		t.Parallel()
		t.Run("bad request", func(t *testing.T) {
			t.Parallel()
			s := setUp(t)
			defer s.tearDown()
			req, _ := http.NewRequestWithContext(ctx, "GET", "/pickpoint/0", strings.NewReader(""))
			w := httptest.NewRecorder()
			m := mux.NewRouter()
			s.mockServ.EXPECT().Read(gomock.Any(), invalidID).Return(nil, model.ErrorInvalidInput)
			m.HandleFunc("/pickpoint/{point:[0-9]+}", s.handl.Read)
			m.ServeHTTP(w, req)

			s.handl.Read(w, req)

			require.Equal(t, http.StatusBadRequest, w.Code)
			assert.Equal(t, w.Body.String(), "")
		})
		t.Run("not found", func(t *testing.T) {
			t.Parallel()
			s := setUp(t)
			defer s.tearDown()
			req, _ := http.NewRequestWithContext(ctx, "GET", "/pickpoint/100", strings.NewReader(""))
			w := httptest.NewRecorder()
			m := mux.NewRouter()
			s.mockServ.EXPECT().Read(gomock.Any(), id).Return(nil, model.ErrorObjectNotFound)
			m.HandleFunc("/pickpoint/{point:[0-9]+}", s.handl.Read)
			m.ServeHTTP(w, req)

			s.handl.Read(w, req)

			require.Equal(t, http.StatusNotFound, w.Code)
			assert.Equal(t, w.Body.String(), "")
		})
		t.Run("internal error", func(t *testing.T) {
			t.Parallel()
			s := setUp(t)
			defer s.tearDown()
			req, _ := http.NewRequestWithContext(ctx, "GET", "/pickpoint/100", strings.NewReader(""))
			w := httptest.NewRecorder()
			m := mux.NewRouter()
			s.mockServ.EXPECT().Read(gomock.Any(), id).Return(nil, assert.AnError)
			m.HandleFunc("/pickpoint/{point:[0-9]+}", s.handl.Read)
			m.ServeHTTP(w, req)

			s.handl.Read(w, req)

			require.Equal(t, http.StatusInternalServerError, w.Code)
			assert.Equal(t, w.Body.String(), "")
		})
	})
}
