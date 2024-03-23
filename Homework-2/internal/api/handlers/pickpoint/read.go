package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gitlab.ozon.dev/mer_marat/homework/internal/config"
	"gitlab.ozon.dev/mer_marat/homework/internal/model"
)

func Read(ctx context.Context, s service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
		point, err := s.Read(ctx, id)
		if err != nil {
			if errors.Is(err, model.ErrorObjectNotFound) {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			if errors.Is(err, model.ErrorInvalidInput) {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		pointJSON, _ := json.Marshal(point)
		w.Write(pointJSON)
	}
}
