package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"homework2/pup/internal/model"
	"homework2/pup/internal/pkg/db"
	"homework2/pup/internal/pkg/repository/postgres"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

const (
	queryParamKey = "point"
	port          = ":9000"
)

type server struct {
	repo *postgres.PickpointRepo
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	database, err := db.NewDB(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	repo := postgres.NewRepo(database)
	serv := server{repo: repo}

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
			id, err := serv.repo.Add(ctx, pointRepo)
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
			tag, err := serv.repo.Update(ctx, &point)
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
			point, err := serv.repo.GetByID(ctx, id)
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
			tag, err := serv.repo.Delete(ctx, id)
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

	http.Handle("/", router)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}
