package server

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"gitlab.ozon.dev/mer_marat/homework/cmd/config"
	"gitlab.ozon.dev/mer_marat/homework/internal/model"
	"gitlab.ozon.dev/mer_marat/homework/internal/pkg/repository/postgres"

	"github.com/gorilla/mux"
)

type Server struct {
	repo *postgres.PickpointRepo
}

func NewServer(repo *postgres.PickpointRepo) Server {
	return Server{repo: repo}
}

func (s Server) Create(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var point model.PickPointAdd
	if err = json.Unmarshal(body, &point); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	pointRepo := &model.PickPoint{
		Name:    point.Name,
		Address: point.Address,
		Contact: point.Contact,
	}
	id, err := s.repo.Add(ctx, pointRepo)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	pointRepo.ID = id
	pointJSON, _ := json.Marshal(pointRepo)
	w.Write(pointJSON)
}

func (s Server) Update(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var point model.PickPoint
	if err = json.Unmarshal(body, &point); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	tag, err := s.repo.Update(ctx, &point)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(tag)
}

func (s Server) Read(ctx context.Context, w http.ResponseWriter, r *http.Request) {
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
	point, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, model.ErrorObjectNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	pointJSON, _ := json.Marshal(point)
	w.Write(pointJSON)
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
