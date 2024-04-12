package router

import (
	"fmt"
	"net/http"

	"gitlab.ozon.dev/mer_marat/homework/internal/config"

	"github.com/gorilla/mux"
)

type handler interface {
	Create(w http.ResponseWriter, r *http.Request)
	Read(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type middleware interface {
	AuthMiddleWare(handler http.Handler) http.Handler
	LogMiddleWare(handler http.Handler) http.Handler
}

func MakeRouter(h handler, middleware middleware, cfg config.Config) *mux.Router {
	router := mux.NewRouter()
	router.Use(middleware.LogMiddleWare)
	router.Use(middleware.AuthMiddleWare)
	router.HandleFunc("/pickpoint", h.Create).Methods("POST")
	router.HandleFunc("/pickpoint", h.Update).Methods("PUT")
	router.HandleFunc(fmt.Sprintf("/pickpoint/{%s:[0-9]+}", config.QueryParamKey), h.Delete).Methods("DELETE")
	router.HandleFunc(fmt.Sprintf("/pickpoint/{%s:[0-9]+}", config.QueryParamKey), h.Read).Methods("GET")
	return router
}
