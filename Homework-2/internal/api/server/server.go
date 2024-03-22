package server

import (
	"context"
	"net/http"
	"strconv"

	"gitlab.ozon.dev/mer_marat/homework/internal/api/router"
	"gitlab.ozon.dev/mer_marat/homework/internal/config"
	"gitlab.ozon.dev/mer_marat/homework/internal/service/pickpoint"
	"gitlab.ozon.dev/mer_marat/homework/internal/storage/postgres"

	"github.com/gorilla/mux"
)

type Server struct {
	repo *postgres.PickpointRepo
}

func NewServer(repo *postgres.PickpointRepo) Server {
	return Server{repo: repo}
}

func (s Server) Run(ctx context.Context, cfg config.Config, cancel context.CancelFunc) error {
	defer cancel()
	serv := pickpoint.NewServiceRepo(s.repo)
	router := router.MakeRouter(ctx, serv, cfg)

	http.Handle("/", router)
	return http.ListenAndServe(cfg.Server.Port, nil)
}

func (s Server) Delete(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key, ok := vars[config.QueryParamKey]
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
	}
	id, err := strconv.ParseInt(key, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	tag, err := s.repo.Delete(ctx, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(tag)
}
