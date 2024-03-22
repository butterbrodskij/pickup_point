package pickpointhandler

import (
	"context"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"gitlab.ozon.dev/mer_marat/homework/internal/config"
	"gitlab.ozon.dev/mer_marat/homework/internal/model"
	"gitlab.ozon.dev/mer_marat/homework/internal/service/pickpoint"
)

func Delete(ctx context.Context, s pickpoint.ServiceRepo, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key, ok := vars[config.QueryParamKey]
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
	}
	tag, err := s.Delete(ctx, key)
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
