package handler

import (
	"context"
	"log"
	"net/http"

	pickpointhandler "gitlab.ozon.dev/mer_marat/homework/internal/api/handler/pickpoint_handler"
	"gitlab.ozon.dev/mer_marat/homework/internal/config"
	"gitlab.ozon.dev/mer_marat/homework/internal/service/pickpoint"
)

func AuthMiddleWare(handler http.Handler, cfg config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		login, password, ok := r.BasicAuth()
		if !ok || !isValidUser(login, password, cfg) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		handler.ServeHTTP(w, r)
	}
}

func LogMiddleWare(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("got request with method %s\n", r.Method)
		handler.ServeHTTP(w, r)
	})
}

func PickpointHandler(ctx context.Context, serv pickpoint.ServiceRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			pickpointhandler.Create(ctx, serv, w, r)
		case http.MethodPut:
			pickpointhandler.Update(ctx, serv, w, r)
		case http.MethodGet:
			fallthrough
		case http.MethodDelete:
			fallthrough
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func PickpointKeyHandler(ctx context.Context, serv pickpoint.ServiceRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			pickpointhandler.Read(ctx, serv, w, r)
		case http.MethodDelete:
			pickpointhandler.Delete(ctx, serv, w, r)
		case http.MethodPost:
			fallthrough
		case http.MethodPut:
			fallthrough
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func isValidUser(login, password string, cfg config.Config) bool {
	for _, user := range cfg.Users {
		if user.Login == login && user.Password == password {
			return true
		}
	}
	return false
}
