//go:generate mockgen -source=./sender.go -destination=./sender_mocks_test.go -package=kafka
package kafka

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/IBM/sarama"
	"gitlab.ozon.dev/mer_marat/homework/internal/model"
)

type producer interface {
	SendSyncMessage(message *sarama.ProducerMessage) (partition int32, offset int64, err error)
}

type KafkaSender struct {
	producer producer
	topic    string
}

func NewKafkaSender(producer producer, topic string) *KafkaSender {
	return &KafkaSender{
		producer: producer,
		topic:    topic,
	}
}

func (s *KafkaSender) SendMessage(message model.RequestMessage) error {
	kafkaMsg, err := s.buildMessage(message)
	if err != nil {
		return err
	}
	if _, _, err = s.producer.SendSyncMessage(kafkaMsg); err != nil {
		return err
	}
	return nil
}

func (s *KafkaSender) buildMessage(message model.RequestMessage) (*sarama.ProducerMessage, error) {
	kafkaMsg, err := convert2KafkaMessage(message)
	if err != nil {
		return nil, err
	}
	msg, err := json.Marshal(kafkaMsg)
	if err != nil {
		return nil, err
	}
	return &sarama.ProducerMessage{
		Topic:     s.topic,
		Value:     sarama.ByteEncoder(msg),
		Partition: -1,
		Key:       sarama.StringEncoder(message.Request.Method),
		Headers: []sarama.RecordHeader{
			{
				Key:   []byte("test-header"),
				Value: []byte("test-value"),
			},
		},
	}, nil
}

func convert2KafkaMessage(msg model.RequestMessage) (model.LogMessage, error) {
	if msg.Request == nil {
		return model.LogMessage{}, model.ErrorEmptyRequest
	}
	if msg.Request.Body == nil {
		return model.LogMessage{}, model.ErrorEmptyBodyRequest
	}
	reqBytes, err := io.ReadAll(msg.Request.Body)
	if err != nil {
		return model.LogMessage{}, err
	}
	reqString := string(reqBytes)

	login, _, _ := msg.Request.BasicAuth()
	msg.Request.Body = io.NopCloser(bytes.NewBuffer(reqBytes))

	return model.LogMessage{
		CaughtTime: msg.CaughtTime,
		Method:     msg.Request.Method,
		Url:        msg.Request.URL,
		Body:       reqString,
		Login:      login,
	}, nil
}
