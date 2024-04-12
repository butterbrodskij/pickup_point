package kafka

import (
	"github.com/IBM/sarama"
)

type ConsumerGroup struct {
	ready    chan bool
	handlers map[string]Handler
	name     string
}

func NewConsumerGroup(handlers map[string]Handler, groupName string) *ConsumerGroup {
	return &ConsumerGroup{
		ready:    make(chan bool),
		handlers: handlers,
		name:     groupName,
	}
}

func (consumer *ConsumerGroup) Ready() <-chan bool {
	return consumer.ready
}

func (consumer *ConsumerGroup) Setup(_ sarama.ConsumerGroupSession) error {
	close(consumer.ready)

	return nil
}

func (consumer *ConsumerGroup) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (consumer *ConsumerGroup) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message := <-claim.Messages():
			if handler, ok := consumer.handlers[message.Topic]; ok {
				handler.Handle(message)
			}
			session.MarkMessage(message, "")
		case <-session.Context().Done():
			return nil
		}
	}
}
