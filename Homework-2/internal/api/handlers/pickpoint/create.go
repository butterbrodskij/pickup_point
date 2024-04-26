package handler

import (
	"encoding/json"
	"net/http"

	"gitlab.ozon.dev/mer_marat/homework/internal/model"
	pickpoint_pb "gitlab.ozon.dev/mer_marat/homework/internal/pkg/pb/pickpoint"
)

func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	var point model.PickPoint
	if err := json.NewDecoder(r.Body).Decode(&point); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	pointNew, err := h.service.Create(r.Context(), &pickpoint_pb.PickPoint{
		Id:      point.ID,
		Name:    point.Name,
		Address: point.Address,
		Contact: point.Contact,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	pointJSON, _ := json.Marshal(pointNew)
	w.Write(pointJSON)
}
