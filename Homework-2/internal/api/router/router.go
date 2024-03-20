package router

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"gitlab.ozon.dev/mer_marat/homework/internal/api/server"
	"gitlab.ozon.dev/mer_marat/homework/internal/model"

	"github.com/gorilla/mux"
)

const queryParamKey = "point"

func MakeRouter(ctx context.Context, serv server.Server) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/pickpoint", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
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
			pointRepo := &model.PickPoint{
				Name:    point.Name,
				Address: point.Address,
				Contact: point.Contact,
			}
			id, err := serv.Create(ctx, pointRepo)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			pointRepo.ID = id
			pointJSON, _ := json.Marshal(pointRepo)
			w.Write(pointJSON)
		case http.MethodPut:
			body, err := io.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			var point model.PickPoint
			if err = json.Unmarshal(body, &point); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			tag, err := serv.Update(ctx, &point)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Write(tag)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
	})

	router.HandleFunc(fmt.Sprintf("/pickpoint/{%s:[0-9]+}", queryParamKey), func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			vars := mux.Vars(r)
			id, err := strconv.ParseInt(vars[queryParamKey], 10, 64)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			point, err := serv.Read(ctx, id)
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
		case http.MethodDelete:
			vars := mux.Vars(r)
			id, err := strconv.ParseInt(vars[queryParamKey], 10, 64)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			tag, err := serv.Delete(ctx, id)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Write(tag)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
	})

	return router
}
