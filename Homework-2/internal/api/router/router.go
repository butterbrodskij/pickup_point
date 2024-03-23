package router

import (
	"fmt"
	"net/http"

	handler "gitlab.ozon.dev/mer_marat/homework/internal/api/handlers/pickpoint"
	"gitlab.ozon.dev/mer_marat/homework/internal/api/middleware"
	"gitlab.ozon.dev/mer_marat/homework/internal/config"
	"gitlab.ozon.dev/mer_marat/homework/internal/service/pickpoint"

	"github.com/gorilla/mux"
)

func MakeRouter(serv pickpoint.Service, cfg config.Config) *mux.Router {
	router := mux.NewRouter()
	router.Use(middleware.LogMiddleWare)
	router.Use(func(h http.Handler) http.Handler {
		return middleware.AuthMiddleWare(h, cfg)
	})
	router.HandleFunc("/pickpoint", handler.Create(serv)).Methods("POST")
	router.HandleFunc("/pickpoint", handler.Update(serv)).Methods("PUT")
	router.HandleFunc(fmt.Sprintf("/pickpoint/{%s:[0-9]+}", config.QueryParamKey), handler.Delete(serv)).Methods("DELETE")
	router.HandleFunc(fmt.Sprintf("/pickpoint/{%s:[0-9]+}", config.QueryParamKey), handler.Read(serv)).Methods("GET")
	return router
}
