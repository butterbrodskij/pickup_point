package handler

import (
	"context"
	"net/http"

	pickpointhandler "gitlab.ozon.dev/mer_marat/homework/internal/api/handler/pickpoint_handler"
	"gitlab.ozon.dev/mer_marat/homework/internal/service/pickpoint"
)

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
