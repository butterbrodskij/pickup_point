package middleware

import (
	"log"
	"net/http"

	"gitlab.ozon.dev/mer_marat/homework/internal/config"
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

func isValidUser(login, password string, cfg config.Config) bool {
	for _, user := range cfg.Users {
		if user.Login == login && user.Password == password {
			return true
		}
	}
	return false
}
