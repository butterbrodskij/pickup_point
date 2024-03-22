package pickpointhandler

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"gitlab.ozon.dev/mer_marat/homework/internal/model"
	"gitlab.ozon.dev/mer_marat/homework/internal/service/pickpoint"
)

func Create(ctx context.Context, s pickpoint.ServiceRepo, w http.ResponseWriter, r *http.Request) {
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
	pointRepo, err := s.Create(ctx, point)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	pointJSON, _ := json.Marshal(pointRepo)
	w.Write(pointJSON)
}
