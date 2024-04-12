package dummy

import (
	"encoding/json"
	"fmt"

	"github.com/IBM/sarama"
	"gitlab.ozon.dev/mer_marat/homework/internal/model"
)

type Handler struct {
	wait   chan struct{}
	Result string
	Err    error
}

func NewHandler() *Handler {
	return &Handler{
		wait: make(chan struct{}, 1),
	}
}

func (l *Handler) Wait() <-chan struct{} {
	return l.wait
}

func (l *Handler) Handle(message *sarama.ConsumerMessage) {
	logMsg, err := consumerMessage2LogMessage(message)
	if err != nil {
		l.Err = err
		l.Result = ""
		return
	}
	l.Err = nil
	l.Result = fmt.Sprintf(`New Request:
	Method: %s	Path: %s	login: %s	Body: %s`,
		logMsg.Method, logMsg.Url.Path, logMsg.Login, logMsg.Body)
	l.wait <- struct{}{}
}

func consumerMessage2LogMessage(message *sarama.ConsumerMessage) (model.LogMessage, error) {
	var logMsg model.LogMessage
	err := json.Unmarshal(message.Value, &logMsg)
	return logMsg, err
}
