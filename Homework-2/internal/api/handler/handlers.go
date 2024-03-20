package handler

import (
	"context"
	"log"
	"net/http"

	"gitlab.ozon.dev/mer_marat/homework/internal/api/server"
)

func LogMiddleWare(handler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("got request with method %s\n", r.Method)
		handler.ServeHTTP(w, r)
	}
}

func PickpointHandler(ctx context.Context, serv server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			serv.Create(ctx, w, r)
		case http.MethodPut:
			serv.Update(ctx, w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func PickpointKeyHandler(ctx context.Context, serv server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			serv.Read(ctx, w, r)
		case http.MethodDelete:
			serv.Delete(ctx, w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}
