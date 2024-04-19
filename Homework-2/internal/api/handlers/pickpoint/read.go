package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gitlab.ozon.dev/mer_marat/homework/internal/config"
	"gitlab.ozon.dev/mer_marat/homework/internal/model"
)

func (h *handler) Read(w http.ResponseWriter, r *http.Request) {
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
	body, err := h.cache.Get(r.Context(), key)
	if err == nil {
		w.Write([]byte(body))
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
	err = h.cache.Set(r.Context(), key, string(pointJSON))
	if err != nil {
		log.Printf("cache set failed: %s", err)
	}
	w.Write(pointJSON)
}
