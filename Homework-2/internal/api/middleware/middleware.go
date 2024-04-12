package middleware

import (
	"log"
	"net/http"
	"time"

	"gitlab.ozon.dev/mer_marat/homework/internal/config"
	"gitlab.ozon.dev/mer_marat/homework/internal/model"
)

type sender interface {
	SendMessage(message model.RequestMessage) error
}

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

func LogMiddleWare(handler http.Handler, sender sender) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := sender.SendMessage(model.RequestMessage{
			CaughtTime: time.Now(),
			Request:    r,
		})
		if err != nil {
			log.Println(err)
			return
		}
		handler.ServeHTTP(w, r)
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
