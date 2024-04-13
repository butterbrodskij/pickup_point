package kafka

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/IBM/sarama"
)

type KafkaGroupReceiver struct {
	consumer *ConsumerGroup
	client   sarama.ConsumerGroup
	wg       *sync.WaitGroup
	ctx      context.Context
	cancel   context.CancelFunc
}

func NewReceiverGroup(ctxParent context.Context, consumer *ConsumerGroup, brokers []string) (*KafkaGroupReceiver, error) {
	config := sarama.NewConfig()
	config.Version = sarama.MaxVersion

	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	config.Consumer.Group.ResetInvalidOffsets = true

	config.Consumer.Group.Heartbeat.Interval = 3 * time.Second

	config.Consumer.Group.Session.Timeout = 60 * time.Second

	config.Consumer.Group.Rebalance.Timeout = 60 * time.Second

	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}

	client, err := sarama.NewConsumerGroup(brokers, consumer.name, config)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithCancel(ctxParent)

	return &KafkaGroupReceiver{
		consumer: consumer,
		client:   client,
		wg:       &sync.WaitGroup{},
		ctx:      ctx,
		cancel:   cancel,
	}, nil
}

func (r *KafkaGroupReceiver) Close() error {
	r.cancel()
	r.wg.Wait()

	return r.client.Close()
}

func (r *KafkaGroupReceiver) Subscribe(topics []string) error {
	r.wg.Add(1)
	go func() {
		defer r.wg.Done()
		for {
			if err := r.client.Consume(r.ctx, topics, r.consumer); err != nil {
				log.Printf("Error from consumer: %v", err)
			}
			if r.ctx.Err() != nil {
				return
			}
		}
	}()

	<-r.consumer.Ready()
	return nil
}
