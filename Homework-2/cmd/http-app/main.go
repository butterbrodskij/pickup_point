package main

import (
	"context"
	"log"

	"gitlab.ozon.dev/mer_marat/homework/internal/api/server"
	"gitlab.ozon.dev/mer_marat/homework/internal/config"
	"gitlab.ozon.dev/mer_marat/homework/internal/infrastructure/kafka"
	"gitlab.ozon.dev/mer_marat/homework/internal/pkg/db"
	"gitlab.ozon.dev/mer_marat/homework/internal/service/logger"
	"gitlab.ozon.dev/mer_marat/homework/internal/service/pickpoint"
	"gitlab.ozon.dev/mer_marat/homework/internal/storage/postgres"
)

func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(cfg)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	database, err := db.NewDB(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	repo := postgres.NewRepo(database)
	service := pickpoint.NewService(repo)
	serv := server.NewServer(service)

	consumer, err := kafka.NewConsumer(cfg.Brokers)
	if err == nil {
		logger := logger.NewLogger()
		receiver := kafka.NewReceiver(consumer, logger)
		err = receiver.Subscribe("logs")
		if err != nil {
			log.Println("Kafka Subscription Failed")
			cfg.Brokers = nil
		}
	}

	if err := serv.Run(ctx, cfg); err != nil {
		log.Fatal(err)
	}
}
