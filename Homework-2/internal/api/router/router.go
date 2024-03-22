package router

import (
	"context"
	"fmt"
	"net/http"

	"gitlab.ozon.dev/mer_marat/homework/internal/api/handler"
	"gitlab.ozon.dev/mer_marat/homework/internal/api/middleware"
	"gitlab.ozon.dev/mer_marat/homework/internal/config"
	"gitlab.ozon.dev/mer_marat/homework/internal/service/pickpoint"

	"github.com/gorilla/mux"
)

func MakeRouter(ctx context.Context, serv pickpoint.ServiceRepo, cfg config.Config) *mux.Router {
	router := mux.NewRouter()
	router.Use(middleware.LogMiddleWare)
	router.Use(func(h http.Handler) http.Handler {
		return middleware.AuthMiddleWare(h, cfg)
	})
	router.HandleFunc("/pickpoint", handler.PickpointHandler(ctx, serv))
	router.HandleFunc(fmt.Sprintf("/pickpoint/{%s:[0-9]+}", config.QueryParamKey), handler.PickpointKeyHandler(ctx, serv))
	return router
}
