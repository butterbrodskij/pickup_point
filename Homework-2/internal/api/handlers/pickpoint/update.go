package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"gitlab.ozon.dev/mer_marat/homework/internal/model"
	pickpoint_pb "gitlab.ozon.dev/mer_marat/homework/internal/pkg/pb/pickpoint"
)

func (h *handler) Update(w http.ResponseWriter, r *http.Request) {
	var point model.PickPoint
	if err := json.NewDecoder(r.Body).Decode(&point); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, err := h.service.Update(r.Context(), &pickpoint_pb.PickPoint{
		Id:      point.ID,
		Name:    point.Name,
		Address: point.Address,
		Contact: point.Contact,
	})
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
