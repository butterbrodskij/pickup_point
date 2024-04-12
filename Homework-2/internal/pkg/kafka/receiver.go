package kafka

import (
	"github.com/IBM/sarama"
)

type Handler interface {
	Handle(message *sarama.ConsumerMessage)
}

type KafkaReceiver struct {
	consumer *Consumer
	Handler
}

func NewReceiver(consumer *Consumer, handler Handler) *KafkaReceiver {
	return &KafkaReceiver{
		consumer: consumer,
		Handler:  handler,
	}
}

func (r *KafkaReceiver) Subscribe(topic string) error {
	partitionList, err := r.consumer.Consumer.Partitions(topic)

	if err != nil {
		return err
	}

	initialOffset := sarama.OffsetNewest

	for _, partition := range partitionList {
		pc, err := r.consumer.Consumer.ConsumePartition(topic, partition, initialOffset)

		if err != nil {
			return err
		}

		go func(pc sarama.PartitionConsumer) {
			for message := range pc.Messages() {
				r.Handle(message)
			}
		}(pc)
	}

	return nil
}
