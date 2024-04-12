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

type Middleware struct {
	cfg config.Config
	sender
}

func NewMiddleware(cfg config.Config, sender sender) *Middleware {
	return &Middleware{
		cfg:    cfg,
		sender: sender,
	}
}

func (m *Middleware) AuthMiddleWare(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		login, password, ok := r.BasicAuth()
		if !ok || !isValidUser(login, password, m.cfg) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		handler.ServeHTTP(w, r)
	})
}

func (m *Middleware) LogMiddleWare(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := m.SendMessage(model.RequestMessage{
			CaughtTime: time.Now(),
			Request:    r,
		})
		if err != nil {
			log.Println(err)
			return
		}
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
