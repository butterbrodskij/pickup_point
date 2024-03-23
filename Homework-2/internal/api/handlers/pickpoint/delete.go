package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gitlab.ozon.dev/mer_marat/homework/internal/config"
	"gitlab.ozon.dev/mer_marat/homework/internal/model"
)

func Delete(s service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
		err = s.Delete(r.Context(), id)
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
}
