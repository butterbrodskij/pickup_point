package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"gitlab.ozon.dev/mer_marat/homework/internal/model"
)

func (h handler) Update(w http.ResponseWriter, r *http.Request) {
	var point model.PickPoint
	if err := json.NewDecoder(r.Body).Decode(&point); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err := h.service.Update(r.Context(), &point)
	if err != nil {
		if errors.Is(err, model.ErrorInvalidInput) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if errors.Is(err, model.ErrorObjectNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(model.MessageSuccess)
}
