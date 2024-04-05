//go:build integration
// +build integration

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
	"gitlab.ozon.dev/mer_marat/homework/internal/model"
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
		defer db.TearDown(t, "pickpoints")
		handl := handler.NewHandler(pickpoint.NewService(postgres.NewRepo(db.DB)))
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
			defer db.TearDown(t, "pickpoints")
			handl := handler.NewHandler(pickpoint.NewService(postgres.NewRepo(db.DB)))
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
	t.Run("successful deleting pickpoint", func(t *testing.T) {
		db.SetUp(t, "pickpoints")
		defer db.TearDown(t, "pickpoints")
		handl := handler.NewHandler(pickpoint.NewService(postgres.NewRepo(db.DB)))
		fillDB(fixture.PickPoint().Valid1().P())
		req, _ := http.NewRequestWithContext(ctx, "DELETE", "/pickpoint/100", strings.NewReader(""))
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
			defer db.TearDown(t, "pickpoints")
			handl := handler.NewHandler(pickpoint.NewService(postgres.NewRepo(db.DB)))
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
			defer db.TearDown(t, "pickpoints")
			handl := handler.NewHandler(pickpoint.NewService(postgres.NewRepo(db.DB)))
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

func TestRead(t *testing.T) {
	var (
		ctx = context.Background()
	)
	t.Run("successful reading pickpoint", func(t *testing.T) {
		db.SetUp(t, "pickpoints")
		defer db.TearDown(t, "pickpoints")
		handl := handler.NewHandler(pickpoint.NewService(postgres.NewRepo(db.DB)))
		fillDB(fixture.PickPoint().Valid1().P())
		req, _ := http.NewRequestWithContext(ctx, "GET", "/pickpoint/100", strings.NewReader(""))
		w := httptest.NewRecorder()
		m := mux.NewRouter()
		m.HandleFunc("/pickpoint/{point:[0-9]+}", handl.Read)
		m.ServeHTTP(w, req)

		handl.Read(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, w.Body.String(), `{"id":100,"name":"Chertanovo","address":"Chertanovskaya street, 8","contacts":"+7(999)888-77-66"}`)
	})
	t.Run("fail", func(t *testing.T) {
		t.Run("bad request", func(t *testing.T) {
			db.SetUp(t, "pickpoints")
			defer db.TearDown(t, "pickpoints")
			handl := handler.NewHandler(pickpoint.NewService(postgres.NewRepo(db.DB)))
			req, _ := http.NewRequestWithContext(ctx, "GET", "/pickpoint/0", strings.NewReader(""))
			w := httptest.NewRecorder()
			m := mux.NewRouter()
			m.HandleFunc("/pickpoint/{point:[0-9]+}", handl.Read)
			m.ServeHTTP(w, req)

			handl.Read(w, req)

			require.Equal(t, w.Code, http.StatusBadRequest)
			assert.Equal(t, w.Body.String(), "")
		})
		t.Run("not found", func(t *testing.T) {
			db.SetUp(t, "pickpoints")
			defer db.TearDown(t, "pickpoints")
			handl := handler.NewHandler(pickpoint.NewService(postgres.NewRepo(db.DB)))
			req, _ := http.NewRequestWithContext(ctx, "GET", "/pickpoint/1", strings.NewReader(""))
			w := httptest.NewRecorder()
			m := mux.NewRouter()
			m.HandleFunc("/pickpoint/{point:[0-9]+}", handl.Read)
			m.ServeHTTP(w, req)

			handl.Read(w, req)

			require.Equal(t, w.Code, http.StatusNotFound)
			assert.Equal(t, w.Body.String(), "")
		})
	})
}

func TestUpdate(t *testing.T) {
	var (
		ctx = context.Background()
	)
	t.Run("successful updating pickpoint", func(t *testing.T) {
		db.SetUp(t, "pickpoints")
		defer db.TearDown(t, "pickpoints")
		handl := handler.NewHandler(pickpoint.NewService(postgres.NewRepo(db.DB)))
		fillDB(fixture.PickPoint().Valid1().P())
		body := `{"id":100, "name":"Chertanovo", "address":"Chertanovskaya street, 13", "contacts":"+7(999)888-77-66"}`
		req, _ := http.NewRequestWithContext(ctx, "PUT", "/pickpoint", strings.NewReader(body))
		w := httptest.NewRecorder()

		handl.Update(w, req)

		require.Equal(t, w.Code, http.StatusOK)
		assert.Equal(t, w.Body.String(), `operation completed successfully`)
	})
	t.Run("fail", func(t *testing.T) {
		t.Run("bad request", func(t *testing.T) {
			db.SetUp(t, "pickpoints")
			defer db.TearDown(t, "pickpoints")
			handl := handler.NewHandler(pickpoint.NewService(postgres.NewRepo(db.DB)))
			body := `example of bad request`
			req, _ := http.NewRequestWithContext(ctx, "PUT", "/pickpoint", strings.NewReader(body))
			w := httptest.NewRecorder()

			handl.Update(w, req)

			require.Equal(t, w.Code, http.StatusBadRequest)
			assert.Equal(t, w.Body.String(), "")
		})
		t.Run("not found", func(t *testing.T) {
			db.SetUp(t, "pickpoints")
			defer db.TearDown(t, "pickpoints")
			handl := handler.NewHandler(pickpoint.NewService(postgres.NewRepo(db.DB)))
			body := `{"id":100, "name":"Chertanovo", "address":"Chertanovskaya street, 13", "contacts":"+7(999)888-77-66"}`
			req, _ := http.NewRequestWithContext(ctx, "PUT", "/pickpoint", strings.NewReader(body))
			w := httptest.NewRecorder()

			handl.Update(w, req)

			require.Equal(t, w.Code, http.StatusNotFound)
			assert.Equal(t, w.Body.String(), "")
		})
	})
}

func fillDB(point *model.PickPoint) {
	db.DB.ExecQueryRow(context.Background(), "INSERT INTO pickpoints(id, name, address, contacts) VALUES ($1, $2, $3, $4)", point.ID, point.Name, point.Address, point.Contact)
}
