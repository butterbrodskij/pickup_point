package main

import (
	"context"
	"log"

	"gitlab.ozon.dev/mer_marat/homework/internal/api/server"
	"gitlab.ozon.dev/mer_marat/homework/internal/config"
	"gitlab.ozon.dev/mer_marat/homework/internal/pkg/db"
	inmemorycache "gitlab.ozon.dev/mer_marat/homework/internal/pkg/in_memory_cache"
	"gitlab.ozon.dev/mer_marat/homework/internal/pkg/kafka"
	"gitlab.ozon.dev/mer_marat/homework/internal/pkg/redis"
	"gitlab.ozon.dev/mer_marat/homework/internal/service/logger"
	"gitlab.ozon.dev/mer_marat/homework/internal/service/pickpoint"
	"gitlab.ozon.dev/mer_marat/homework/internal/storage/postgres"
)

func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	database, err := db.NewDB(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	redis := redis.NewRedisDB(cfg)
	if err := redis.Ping(ctx); err != nil {
		log.Fatal(err)
	}

	repo := postgres.NewRepo(database)
	cache := inmemorycache.NewInMemoryCache()
	defer cache.Close()
	//service := pickpoint.NewService(repo, redis, database)
	service := pickpoint.NewService(repo, cache, database)

	handler := logger.NewHandler()
	consumer := kafka.NewConsumerGroup(map[string]kafka.Handler{cfg.Kafka.Topic: handler}, cfg.Kafka.Topic)
	receiver, err := kafka.NewReceiverGroup(ctx, consumer, cfg.Kafka.Brokers)
	if err != nil {
		log.Fatal(err)
	}
	err = receiver.Subscribe([]string{cfg.Kafka.Topic})
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := receiver.Close(); err != nil {
			log.Println(err)
		}
	}()

	producer, err := kafka.NewProducer(cfg.Kafka.Brokers)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := producer.Close(); err != nil {
			log.Println(err)
		}
	}()

	serv := server.NewServer(service, producer)
	log.Println("Ready to run")

	if err := serv.Run(ctx, cfg); err != nil {
		log.Fatal(err)
	}
}
