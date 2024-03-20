package router

import (
	"context"
	"fmt"
	"net/http"

	"gitlab.ozon.dev/mer_marat/homework/cmd/config"
	"gitlab.ozon.dev/mer_marat/homework/internal/api/handler"
	"gitlab.ozon.dev/mer_marat/homework/internal/api/server"

	"github.com/gorilla/mux"
)

func MakeRouter(ctx context.Context, serv server.Server, cfg config.Config) *mux.Router {
	router := mux.NewRouter()
	router.Use(handler.LogMiddleWare)
	router.Use(func(h http.Handler) http.Handler {
		return handler.AuthMiddleWare(h, cfg)
	})
	router.HandleFunc("/pickpoint", handler.PickpointHandler(ctx, serv))
	path := fmt.Sprintf("/pickpoint/{%s:[0-9]+}", config.QueryParamKey)
	router.HandleFunc(path, handler.PickpointKeyHandler(ctx, serv))
	return router
}
