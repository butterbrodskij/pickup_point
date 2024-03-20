package router

import (
	"context"
	"fmt"

	"gitlab.ozon.dev/mer_marat/homework/internal/api/handler"
	"gitlab.ozon.dev/mer_marat/homework/internal/api/server"

	"github.com/gorilla/mux"
)

const queryParamKey = "point"

func MakeRouter(ctx context.Context, serv server.Server) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/pickpoint", handler.LogMiddleWare(handler.PickpointHandler(ctx, serv)))
	router.HandleFunc(fmt.Sprintf("/pickpoint/{%s:[0-9]+}", queryParamKey), handler.LogMiddleWare(handler.PickpointKeyHandler(ctx, serv)))
	return router
}
