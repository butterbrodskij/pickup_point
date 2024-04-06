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
)

func TestDelete(t *testing.T) {
	t.Parallel()
	var (
		ctx       = context.Background()
		id        = int64(1)
		invalidID = int64(0)
	)
	t.Run("delete writes success message", func(t *testing.T) {
		t.Parallel()
		s := setUp(t)
		defer s.tearDown()
		req, _ := http.NewRequestWithContext(ctx, "DELETE", "/pickpoint/1", strings.NewReader(""))
		w := httptest.NewRecorder()
		m := mux.NewRouter()
		s.mockServ.EXPECT().Delete(gomock.Any(), id).Return(nil)
		m.HandleFunc("/pickpoint/{point:[0-9]+}", s.handl.Delete)
		m.ServeHTTP(w, req)

		s.handl.Delete(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, w.Body.String(), `operation completed successfully`)
	})
	t.Run("fail", func(t *testing.T) {
		t.Parallel()
		t.Run("bad request", func(t *testing.T) {
			t.Parallel()
			s := setUp(t)
			defer s.tearDown()
			req, _ := http.NewRequestWithContext(ctx, "DELETE", "/pickpoint/0", strings.NewReader(""))
			w := httptest.NewRecorder()
			m := mux.NewRouter()
			s.mockServ.EXPECT().Delete(gomock.Any(), invalidID).Return(model.ErrorInvalidInput)
			m.HandleFunc("/pickpoint/{point:[0-9]+}", s.handl.Delete)
			m.ServeHTTP(w, req)

			s.handl.Delete(w, req)

			require.Equal(t, http.StatusBadRequest, w.Code)
			assert.Equal(t, w.Body.String(), "")
		})
		t.Run("not found", func(t *testing.T) {
			t.Parallel()
			s := setUp(t)
			defer s.tearDown()
			req, _ := http.NewRequestWithContext(ctx, "DELETE", "/pickpoint/1", strings.NewReader(""))
			w := httptest.NewRecorder()
			m := mux.NewRouter()
			s.mockServ.EXPECT().Delete(gomock.Any(), id).Return(model.ErrorObjectNotFound)
			m.HandleFunc("/pickpoint/{point:[0-9]+}", s.handl.Delete)
			m.ServeHTTP(w, req)

			s.handl.Delete(w, req)

			require.Equal(t, http.StatusNotFound, w.Code)
			assert.Equal(t, w.Body.String(), "")
		})
		t.Run("internal error", func(t *testing.T) {
			t.Parallel()
			s := setUp(t)
			defer s.tearDown()
			req, _ := http.NewRequestWithContext(ctx, "DELETE", "/pickpoint/1", strings.NewReader(""))
			w := httptest.NewRecorder()
			m := mux.NewRouter()
			s.mockServ.EXPECT().Delete(gomock.Any(), id).Return(assert.AnError)
			m.HandleFunc("/pickpoint/{point:[0-9]+}", s.handl.Delete)
			m.ServeHTTP(w, req)

			s.handl.Delete(w, req)

			require.Equal(t, http.StatusInternalServerError, w.Code)
			assert.Equal(t, w.Body.String(), "")
		})
	})
}
