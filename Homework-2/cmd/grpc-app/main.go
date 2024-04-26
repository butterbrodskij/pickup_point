package main

import (
	"context"
	"log"

	"github.com/prometheus/client_golang/prometheus"
	"gitlab.ozon.dev/mer_marat/homework/internal/api/server"
	"gitlab.ozon.dev/mer_marat/homework/internal/config"
	"gitlab.ozon.dev/mer_marat/homework/internal/pkg/db"
	inmemorycache "gitlab.ozon.dev/mer_marat/homework/internal/pkg/in_memory_cache"
	"gitlab.ozon.dev/mer_marat/homework/internal/pkg/kafka"
	"gitlab.ozon.dev/mer_marat/homework/internal/pkg/redis"
	"gitlab.ozon.dev/mer_marat/homework/internal/service/cover"
	"gitlab.ozon.dev/mer_marat/homework/internal/service/order"
	"gitlab.ozon.dev/mer_marat/homework/internal/service/pickpoint"
	storage "gitlab.ozon.dev/mer_marat/homework/internal/storage/file"
	"gitlab.ozon.dev/mer_marat/homework/internal/storage/postgres"
)

var (
	reg = prometheus.NewRegistry()

	pickpointCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "pickpoint_grpc",
		Help: "Number of requests handled",
	})
)

func init() {
	reg.MustRegister(pickpointCounter)
}

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
	service := pickpoint.NewService(repo, cache, database, pickpointCounter)

	producer, err := kafka.NewProducer(cfg.Kafka.Brokers)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := producer.Close(); err != nil {
			log.Println(err)
		}
	}()

	storOrders, err := storage.NewOrders("storage_orders.json")
	if err != nil {
		log.Printf("can not connect to storage: %s\n", err)
		return
	}
	servOrders := order.NewService(&storOrders, cover.NewService())

	serv := server.NewServer(service, servOrders, producer, reg)
	log.Println("Ready to run")

	if err := serv.RunGRPC(ctx, cfg); err != nil {
		log.Fatal(err)
	}
}
