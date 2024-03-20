package router

import (
	"context"
	"fmt"
	"net/http"

	"gitlab.ozon.dev/mer_marat/homework/internal/api/server"

	"github.com/gorilla/mux"
)

const queryParamKey = "point"

func MakeRouter(ctx context.Context, serv server.Server) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/pickpoint", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			serv.Create(ctx, w, r)
		case http.MethodPut:
			serv.Update(ctx, w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	router.HandleFunc(fmt.Sprintf("/pickpoint/{%s:[0-9]+}", queryParamKey), func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			serv.Read(ctx, w, r)
		case http.MethodDelete:
			serv.Delete(ctx, w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	return router
}
