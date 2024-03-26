package handler

import (
	"encoding/json"
	"net/http"

	"gitlab.ozon.dev/mer_marat/homework/internal/model"
)

func Create(s service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var point model.PickPoint
		if err := json.NewDecoder(r.Body).Decode(&point); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		pointNew, err := s.Create(r.Context(), &point)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		pointJSON, _ := json.Marshal(pointNew)
		w.Write(pointJSON)
	}
}
