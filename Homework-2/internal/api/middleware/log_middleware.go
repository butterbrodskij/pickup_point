package middleware

import (
	"log"
	"net/http"
	"time"

	"gitlab.ozon.dev/mer_marat/homework/internal/model"
)

type sender interface {
	SendMessage(message model.RequestMessage) error
}

type LogMiddleware struct {
	sender
}

func NewLogMiddleware(sender sender) *LogMiddleware {
	return &LogMiddleware{
		sender: sender,
	}
}

func (m *LogMiddleware) LogMiddleWare(handler http.Handler) http.Handler {
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
