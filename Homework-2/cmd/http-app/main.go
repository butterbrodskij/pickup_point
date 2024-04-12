package main

import (
	"context"
	"log"

	"gitlab.ozon.dev/mer_marat/homework/internal/api/server"
	"gitlab.ozon.dev/mer_marat/homework/internal/config"
	"gitlab.ozon.dev/mer_marat/homework/internal/pkg/db"
	"gitlab.ozon.dev/mer_marat/homework/internal/pkg/kafka"
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

	repo := postgres.NewRepo(database)
	service := pickpoint.NewService(repo)

	consumer, err := kafka.NewConsumer(cfg.Kafka.Brokers)
	if err != nil {
		log.Fatal(err)
	}

	producer, err := kafka.NewProducer(cfg.Kafka.Brokers)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := producer.Close(); err != nil {
			log.Println(err)
		}
	}()

	serv := server.NewServer(service, consumer, producer)

	if err := serv.Run(ctx, cfg); err != nil {
		log.Fatal(err)
	}
}
