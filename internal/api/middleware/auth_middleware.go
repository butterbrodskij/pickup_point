package middleware

import (
	"net/http"

	"gitlab.ozon.dev/mer_marat/homework/internal/config"
)

type AuthMiddleware struct {
	cfg config.Config
}

func NewAuthMiddleware(cfg config.Config) *AuthMiddleware {
	return &AuthMiddleware{
		cfg: cfg,
	}
}

func (m *AuthMiddleware) AuthMiddleWare(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		login, password, ok := r.BasicAuth()
		if !ok || !isValidUser(login, password, m.cfg) {
			w.WriteHeader(http.StatusUnauthorized)
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
