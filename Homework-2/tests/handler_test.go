//go:build integration

package tests

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	handler "gitlab.ozon.dev/mer_marat/homework/internal/api/handlers/pickpoint"
	"gitlab.ozon.dev/mer_marat/homework/internal/api/middleware"
	"gitlab.ozon.dev/mer_marat/homework/internal/api/router"
	"gitlab.ozon.dev/mer_marat/homework/internal/model"
	"gitlab.ozon.dev/mer_marat/homework/internal/service/pickpoint"
	"gitlab.ozon.dev/mer_marat/homework/internal/storage/postgres"
	"gitlab.ozon.dev/mer_marat/homework/tests/dummy"
	"gitlab.ozon.dev/mer_marat/homework/tests/fixture"
)

func TestCreate(t *testing.T) {
	var (
		ctx            = context.Background()
		authMiddleware = middleware.NewAuthMiddleware(cfg)
		logMiddleware  = middleware.NewLogMiddleware(dummy.NewDummySender())
		service        = pickpoint.NewService(postgres.NewRepo(db.DB), dummy.NewCache(), db.DB)
		router         = router.MakeRouter(handler.NewHandler(service), authMiddleware, logMiddleware, cfg)
	)
	t.Run("creating pickpoint", func(t *testing.T) {
		db.SetUp(t, "pickpoints")
		defer db.TearDown(t, "pickpoints")
		body := `{"name":"Chertanovo", "address":"Chertanovskaya street, 13", "contact":"+7(999)888-77-66"}`
		req, _ := http.NewRequestWithContext(ctx, "POST", "/pickpoint", strings.NewReader(body))
		req.SetBasicAuth("admin", "password")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		require.Equal(t, w.Code, http.StatusOK)
		assert.Equal(t, w.Body.String(), `{"id":1,"name":"Chertanovo","address":"Chertanovskaya street, 13","contact":"+7(999)888-77-66"}`)
	})
	t.Run("fail", func(t *testing.T) {
		t.Run("bad request", func(t *testing.T) {
			db.SetUp(t, "pickpoints")
			defer db.TearDown(t, "pickpoints")
			body := `example of bad request`
			req, _ := http.NewRequestWithContext(ctx, "POST", "/pickpoint", strings.NewReader(body))
			req.SetBasicAuth("admin", "password")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			require.Equal(t, w.Code, http.StatusBadRequest)
			assert.Equal(t, w.Body.String(), "")
		})
		t.Run("unauthorized request", func(t *testing.T) {
			db.SetUp(t, "pickpoints")
			defer db.TearDown(t, "pickpoints")
			body := `{"name":"Chertanovo", "address":"Chertanovskaya street, 13", "contact":"+7(999)888-77-66"}`
			req, _ := http.NewRequestWithContext(ctx, "POST", "/pickpoint", strings.NewReader(body))
			req.SetBasicAuth("admin", "wrong_password")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			require.Equal(t, w.Code, http.StatusUnauthorized)
			assert.Equal(t, w.Body.String(), "")
		})
		t.Run("page not found", func(t *testing.T) {
			db.SetUp(t, "pickpoints")
			defer db.TearDown(t, "pickpoints")
			body := `{"name":"Chertanovo", "address":"Chertanovskaya street, 13", "contact":"+7(999)888-77-66"}`
			req, _ := http.NewRequestWithContext(ctx, "POST", "/wrong_pickpoint", strings.NewReader(body))
			req.SetBasicAuth("admin", "wrong_password")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			require.Equal(t, w.Code, http.StatusNotFound)
			assert.Equal(t, w.Body.String(), "404 page not found\n")
		})
	})
}

func TestDelete(t *testing.T) {
	var (
		ctx            = context.Background()
		authMiddleware = middleware.NewAuthMiddleware(cfg)
		logMiddleware  = middleware.NewLogMiddleware(dummy.NewDummySender())
		service        = pickpoint.NewService(postgres.NewRepo(db.DB), dummy.NewCache(), db.DB)
		router         = router.MakeRouter(handler.NewHandler(service), authMiddleware, logMiddleware, cfg)
	)
	t.Run("successful deleting pickpoint", func(t *testing.T) {
		db.SetUp(t, "pickpoints")
		defer db.TearDown(t, "pickpoints")
		fillDB(fixture.PickPoint().Valid1().P())
		req, _ := http.NewRequestWithContext(ctx, "DELETE", "/pickpoint/100", strings.NewReader(""))
		req.SetBasicAuth("admin", "password")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, w.Body.String(), `operation completed successfully`)
	})
	t.Run("fail", func(t *testing.T) {
		t.Run("bad request", func(t *testing.T) {
			db.SetUp(t, "pickpoints")
			defer db.TearDown(t, "pickpoints")
			req, _ := http.NewRequestWithContext(ctx, "DELETE", "/pickpoint/0", strings.NewReader(""))
			req.SetBasicAuth("admin", "password")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			require.Equal(t, w.Code, http.StatusBadRequest)
			assert.Equal(t, w.Body.String(), "")
		})
		t.Run("not found", func(t *testing.T) {
			db.SetUp(t, "pickpoints")
			defer db.TearDown(t, "pickpoints")
			req, _ := http.NewRequestWithContext(ctx, "DELETE", "/pickpoint/1", strings.NewReader(""))
			req.SetBasicAuth("admin", "password")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			require.Equal(t, w.Code, http.StatusNotFound)
			assert.Equal(t, w.Body.String(), "")
		})
	})
}

func TestRead(t *testing.T) {
	var (
		ctx            = context.Background()
		authMiddleware = middleware.NewAuthMiddleware(cfg)
		logMiddleware  = middleware.NewLogMiddleware(dummy.NewDummySender())
		service        = pickpoint.NewService(postgres.NewRepo(db.DB), dummy.NewCache(), db.DB)
		router         = router.MakeRouter(handler.NewHandler(service), authMiddleware, logMiddleware, cfg)
	)
	t.Run("successful reading pickpoint", func(t *testing.T) {
		db.SetUp(t, "pickpoints")
		defer db.TearDown(t, "pickpoints")
		fillDB(fixture.PickPoint().Valid1().P())
		req, _ := http.NewRequestWithContext(ctx, "GET", "/pickpoint/100", strings.NewReader(""))
		req.SetBasicAuth("admin", "password")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, w.Body.String(), `{"id":100,"name":"Chertanovo","address":"Chertanovskaya street, 8","contact":"+7(999)888-77-66"}`)
	})
	t.Run("fail", func(t *testing.T) {
		t.Run("bad request", func(t *testing.T) {
			db.SetUp(t, "pickpoints")
			defer db.TearDown(t, "pickpoints")
			req, _ := http.NewRequestWithContext(ctx, "GET", "/pickpoint/0", strings.NewReader(""))
			req.SetBasicAuth("admin", "password")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			require.Equal(t, w.Code, http.StatusBadRequest)
			assert.Equal(t, w.Body.String(), "")
		})
		t.Run("not found", func(t *testing.T) {
			db.SetUp(t, "pickpoints")
			defer db.TearDown(t, "pickpoints")
			req, _ := http.NewRequestWithContext(ctx, "GET", "/pickpoint/1", strings.NewReader(""))
			req.SetBasicAuth("admin", "password")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			require.Equal(t, w.Code, http.StatusNotFound)
			assert.Equal(t, w.Body.String(), "")
		})
	})
}

func TestUpdate(t *testing.T) {
	var (
		ctx            = context.Background()
		authMiddleware = middleware.NewAuthMiddleware(cfg)
		logMiddleware  = middleware.NewLogMiddleware(dummy.NewDummySender())
		service        = pickpoint.NewService(postgres.NewRepo(db.DB), dummy.NewCache(), db.DB)
		router         = router.MakeRouter(handler.NewHandler(service), authMiddleware, logMiddleware, cfg)
	)
	t.Run("successful updating pickpoint", func(t *testing.T) {
		db.SetUp(t, "pickpoints")
		defer db.TearDown(t, "pickpoints")
		fillDB(fixture.PickPoint().Valid1().P())
		body := `{"id":100, "name":"Chertanovo", "address":"Chertanovskaya street, 13", "contact":"+7(999)888-77-66"}`
		req, _ := http.NewRequestWithContext(ctx, "PUT", "/pickpoint", strings.NewReader(body))
		req.SetBasicAuth("admin", "password")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		require.Equal(t, w.Code, http.StatusOK)
		assert.Equal(t, w.Body.String(), `operation completed successfully`)
	})
	t.Run("fail", func(t *testing.T) {
		t.Run("bad request", func(t *testing.T) {
			db.SetUp(t, "pickpoints")
			defer db.TearDown(t, "pickpoints")
			body := `example of bad request`
			req, _ := http.NewRequestWithContext(ctx, "PUT", "/pickpoint", strings.NewReader(body))
			req.SetBasicAuth("admin", "password")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			require.Equal(t, w.Code, http.StatusBadRequest)
			assert.Equal(t, w.Body.String(), "")
		})
		t.Run("not found", func(t *testing.T) {
			db.SetUp(t, "pickpoints")
			defer db.TearDown(t, "pickpoints")
			body := `{"id":100, "name":"Chertanovo", "address":"Chertanovskaya street, 13", "contact":"+7(999)888-77-66"}`
			req, _ := http.NewRequestWithContext(ctx, "PUT", "/pickpoint", strings.NewReader(body))
			req.SetBasicAuth("admin", "password")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			require.Equal(t, w.Code, http.StatusNotFound)
			assert.Equal(t, w.Body.String(), "")
		})
	})
}

func fillDB(point *model.PickPoint) {
	db.DB.ExecQueryRow(context.Background(), "INSERT INTO pickpoints(id, name, address, contacts) VALUES ($1, $2, $3, $4)", point.ID, point.Name, point.Address, point.Contact)
}
