package logger

import (
	"context"
	"encoding/json"
	"log"

	"github.com/IBM/sarama"
	"gitlab.ozon.dev/mer_marat/homework/internal/model"
)

type Logger struct {
}

func NewHandler() *Logger {
	return &Logger{}
}

func (l *Logger) Handle(_ context.Context, _ sarama.ConsumerGroupSession, message *sarama.ConsumerMessage) {
	logMsg, err := consumerMessage2LogMessage(message)
	if err != nil {
		log.Println(model.ErrorInvalidKafkaMessage)
		return
	}
	log.Printf(`New Request:
	Caught: %v	Method: %s	Path: %s	login: %s	Body: %s`,
		logMsg.CaughtTime, logMsg.Method, logMsg.Url.Path, logMsg.Login, logMsg.Body)
}

func consumerMessage2LogMessage(message *sarama.ConsumerMessage) (model.LogMessage, error) {
	var logMsg model.LogMessage
	err := json.Unmarshal(message.Value, &logMsg)
	return logMsg, err
}
