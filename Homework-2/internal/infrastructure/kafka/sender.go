package kafka

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/IBM/sarama"
)

type requestMessage struct {
	CaughtTime time.Time
	Request    *http.Request
}

type requestKafkaMessage struct {
	CaughtTime time.Time `json:"time"`
	Method     string    `json:"method"`
	Url        *url.URL  `json:"url"`
	Body       string    `json:"body"`
	Login      string    `json:"login"`
}

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

func (s *KafkaSender) SendMessage(message requestMessage) error {
	kafkaMsg, err := s.buildMessage(message)
	if err != nil {
		return err
	}

	if _, _, err = s.producer.SendSyncMessage(kafkaMsg); err != nil {
		return err
	}

	return nil
}

func (s *KafkaSender) buildMessage(message requestMessage) (*sarama.ProducerMessage, error) {
	msg, err := json.Marshal(convert2KafkaMessage(message))

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

func convert2KafkaMessage(msg requestMessage) requestKafkaMessage {
	reqBytes, _ := io.ReadAll(msg.Request.Body)
	reqString := string(reqBytes)

	login, _, _ := msg.Request.BasicAuth()

	return requestKafkaMessage{
		CaughtTime: msg.CaughtTime,
		Method:     msg.Request.Method,
		Url:        msg.Request.URL,
		Body:       reqString,
		Login:      login,
	}
}
