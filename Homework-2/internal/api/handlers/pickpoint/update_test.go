//go:build !integration

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
	"gitlab.ozon.dev/mer_marat/homework/internal/model"
)

func TestUpdate(t *testing.T) {
	t.Parallel()
	var (
		ctx = context.Background()
	)
	t.Run("smoke test", func(t *testing.T) {
		t.Parallel()
		s := setUp(t)
		defer s.tearDown()
		body := `{"id":100, "name":"Chertanovo", "address":"Chertanovskaya street, 13", "contacts":"+7(999)888-77-66"}`
		req, _ := http.NewRequestWithContext(ctx, "PUT", "/pickpoint", strings.NewReader(body))
		w := httptest.NewRecorder()
		s.mockServ.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)

		s.handl.Update(w, req)

		require.Equal(t, w.Code, http.StatusOK)
		assert.Equal(t, w.Body.String(), "operation completed successfully")
	})
	t.Run("fail", func(t *testing.T) {
		t.Parallel()
		t.Run("bad request", func(t *testing.T) {
			t.Parallel()
			s := setUp(t)
			defer s.tearDown()
			body := `example of bad request`
			req, _ := http.NewRequestWithContext(ctx, "PUT", "/pickpoint", strings.NewReader(body))
			w := httptest.NewRecorder()

			s.handl.Update(w, req)

			require.Equal(t, w.Code, http.StatusBadRequest)
			assert.Equal(t, w.Body.String(), "")
		})
		t.Run("internal error", func(t *testing.T) {
			t.Parallel()
			s := setUp(t)
			defer s.tearDown()
			body := `{"id":100, "name":"Chertanovo", "address":"Chertanovskaya street, 13", "contacts":"+7(999)888-77-66"}`
			req, _ := http.NewRequestWithContext(ctx, "PUT", "/pickpoint", strings.NewReader(body))
			w := httptest.NewRecorder()
			s.mockServ.EXPECT().Update(gomock.Any(), gomock.Any()).Return(assert.AnError)

			s.handl.Update(w, req)

			require.Equal(t, w.Code, http.StatusInternalServerError)
			assert.Equal(t, w.Body.String(), "")
		})
		t.Run("not found", func(t *testing.T) {
			t.Parallel()
			s := setUp(t)
			defer s.tearDown()
			body := `{"id":100, "name":"Chertanovo", "address":"Chertanovskaya street, 13", "contacts":"+7(999)888-77-66"}`
			req, _ := http.NewRequestWithContext(ctx, "PUT", "/pickpoint", strings.NewReader(body))
			w := httptest.NewRecorder()
			s.mockServ.EXPECT().Update(gomock.Any(), gomock.Any()).Return(model.ErrorObjectNotFound)

			s.handl.Update(w, req)

			require.Equal(t, w.Code, http.StatusNotFound)
			assert.Equal(t, w.Body.String(), "")
		})
	})
}
