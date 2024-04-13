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

type authMiddleware interface {
	AuthMiddleWare(handler http.Handler) http.Handler
}

type logMiddleware interface {
	LogMiddleWare(handler http.Handler) http.Handler
}

func MakeRouter(h handler, authMiddleware authMiddleware, logMiddleware logMiddleware, cfg config.Config) *mux.Router {
	router := mux.NewRouter()
	router.Use(logMiddleware.LogMiddleWare)
	router.Use(authMiddleware.AuthMiddleWare)
	router.HandleFunc("/pickpoint", h.Create).Methods("POST")
	router.HandleFunc("/pickpoint", h.Update).Methods("PUT")
	router.HandleFunc(fmt.Sprintf("/pickpoint/{%s:[0-9]+}", config.QueryParamKey), h.Delete).Methods("DELETE")
	router.HandleFunc(fmt.Sprintf("/pickpoint/{%s:[0-9]+}", config.QueryParamKey), h.Read).Methods("GET")
	return router
}
