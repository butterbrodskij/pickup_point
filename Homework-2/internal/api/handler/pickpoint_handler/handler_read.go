package pickpointhandler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"gitlab.ozon.dev/mer_marat/homework/internal/config"
	"gitlab.ozon.dev/mer_marat/homework/internal/model"
	"gitlab.ozon.dev/mer_marat/homework/internal/service/pickpoint"
)

func Read(ctx context.Context, s pickpoint.ServiceRepo, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key, ok := vars[config.QueryParamKey]
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
	}
	point, err := s.Read(ctx, key)
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
