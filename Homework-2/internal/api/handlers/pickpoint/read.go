package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gitlab.ozon.dev/mer_marat/homework/internal/config"
	"gitlab.ozon.dev/mer_marat/homework/internal/model"
)

func (h handler) Read(w http.ResponseWriter, r *http.Request) {
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
	point, err := h.service.Read(r.Context(), id)
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
