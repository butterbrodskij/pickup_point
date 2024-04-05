package tests

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	handler "gitlab.ozon.dev/mer_marat/homework/internal/api/handlers/pickpoint"
	"gitlab.ozon.dev/mer_marat/homework/internal/service/pickpoint"
	"gitlab.ozon.dev/mer_marat/homework/internal/storage/postgres"
	"gitlab.ozon.dev/mer_marat/homework/tests/fixture"
)

func TestCreate(t *testing.T) {
	var (
		ctx = context.Background()
	)
	t.Run("creating pickpoint", func(t *testing.T) {
		db.SetUp(t, "pickpoints")
		defer db.TearDown()
		repo := postgres.NewRepo(db.DB)
		serv := pickpoint.NewService(repo)
		handl := handler.NewHandler(serv)
		body := `{"name":"Chertanovo", "address":"Chertanovskaya street, 13", "contacts":"+7(999)888-77-66"}`
		req, _ := http.NewRequestWithContext(ctx, "POST", "/pickpoint", strings.NewReader(body))
		w := httptest.NewRecorder()

		handl.Create(w, req)

		require.Equal(t, w.Code, http.StatusOK)
		assert.Equal(t, w.Body.String(), `{"id":1,"name":"Chertanovo","address":"Chertanovskaya street, 13","contacts":"+7(999)888-77-66"}`)
	})
	t.Run("fail", func(t *testing.T) {
		t.Run("bad request", func(t *testing.T) {
			db.SetUp(t, "pickpoints")
			defer db.TearDown()
			repo := postgres.NewRepo(db.DB)
			serv := pickpoint.NewService(repo)
			handl := handler.NewHandler(serv)
			body := `example of bad request`
			req, _ := http.NewRequestWithContext(ctx, "POST", "/pickpoint", strings.NewReader(body))
			w := httptest.NewRecorder()

			handl.Create(w, req)

			require.Equal(t, w.Code, http.StatusBadRequest)
			assert.Equal(t, w.Body.String(), "")
		})
	})
}

func TestDelete(t *testing.T) {
	var (
		ctx = context.Background()
	)
	t.Run("deleting pickpoint", func(t *testing.T) {
		db.SetUp(t, "pickpoints")
		defer db.TearDown()
		repo := postgres.NewRepo(db.DB)
		serv := pickpoint.NewService(repo)
		handl := handler.NewHandler(serv)
		_, err := serv.Create(ctx, fixture.PickPoint().ValidInput1().P())
		require.NoError(t, err)
		req, _ := http.NewRequestWithContext(ctx, "DELETE", "/pickpoint/1", strings.NewReader(""))
		w := httptest.NewRecorder()
		m := mux.NewRouter()
		m.HandleFunc("/pickpoint/{point:[0-9]+}", handl.Delete)
		m.ServeHTTP(w, req)

		handl.Delete(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, w.Body.String(), `operation completed successfully`)
	})
	t.Run("fail", func(t *testing.T) {
		t.Run("bad request", func(t *testing.T) {
			db.SetUp(t, "pickpoints")
			defer db.TearDown()
			repo := postgres.NewRepo(db.DB)
			serv := pickpoint.NewService(repo)
			handl := handler.NewHandler(serv)
			req, _ := http.NewRequestWithContext(ctx, "DELETE", "/pickpoint/0", strings.NewReader(""))
			w := httptest.NewRecorder()
			m := mux.NewRouter()
			m.HandleFunc("/pickpoint/{point:[0-9]+}", handl.Delete)
			m.ServeHTTP(w, req)

			handl.Delete(w, req)

			require.Equal(t, w.Code, http.StatusBadRequest)
			assert.Equal(t, w.Body.String(), "")
		})
		t.Run("not found", func(t *testing.T) {
			db.SetUp(t, "pickpoints")
			defer db.TearDown()
			repo := postgres.NewRepo(db.DB)
			serv := pickpoint.NewService(repo)
			handl := handler.NewHandler(serv)
			req, _ := http.NewRequestWithContext(ctx, "DELETE", "/pickpoint/1", strings.NewReader(""))
			w := httptest.NewRecorder()
			m := mux.NewRouter()
			m.HandleFunc("/pickpoint/{point:[0-9]+}", handl.Delete)
			m.ServeHTTP(w, req)

			handl.Delete(w, req)

			require.Equal(t, w.Code, http.StatusNotFound)
			assert.Equal(t, w.Body.String(), "")
		})
	})
}
