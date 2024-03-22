package pickpointhandler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"gitlab.ozon.dev/mer_marat/homework/internal/model"
	"gitlab.ozon.dev/mer_marat/homework/internal/service/pickpoint"
)

func Update(ctx context.Context, s pickpoint.ServiceRepo, w http.ResponseWriter, r *http.Request) {
	var point model.PickPoint
	if err := json.NewDecoder(r.Body).Decode(&point); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	tag, err := s.Update(ctx, &point)
	if err != nil {
		if errors.Is(err, model.ErrorInvalidInput) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(tag)
}
