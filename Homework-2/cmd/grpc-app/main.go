package main

import (
	"context"
	"log"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"gitlab.ozon.dev/mer_marat/homework/internal/api/server"
	"gitlab.ozon.dev/mer_marat/homework/internal/config"
	"gitlab.ozon.dev/mer_marat/homework/internal/pkg/db"
	inmemorycache "gitlab.ozon.dev/mer_marat/homework/internal/pkg/in_memory_cache"
	"gitlab.ozon.dev/mer_marat/homework/internal/pkg/kafka"
	"gitlab.ozon.dev/mer_marat/homework/internal/pkg/redis"
	"gitlab.ozon.dev/mer_marat/homework/internal/pkg/tracer"
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

	givenOrdersCounter = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "given_orders_grpc",
		Help: "Number of given orders",
	})

	requestPickpointMetrics = prometheus.NewHistogram(prometheus.HistogramOpts{
		Namespace: "pickpoint",
		Subsystem: "grpc",
		Name:      "request",
		Help:      "Requests handling histogram",
	})

	failedOrderCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "failed_orders_grpc",
		Help: "Number of failed requests to order service",
	})
)

func init() {
	reg.MustRegister(pickpointCounter, givenOrdersCounter, requestPickpointMetrics, failedOrderCounter)
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
	service := pickpoint.NewService(repo, cache, database)
	service.AddCounterMetric(pickpointCounter)
	service.AddRequestHistogram(requestPickpointMetrics)

	shutdown, err := tracer.InitProvider(ctx, "pickpoint")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := shutdown(ctx); err != nil {
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

	storOrders, err := storage.NewOrders("storage_orders.json")
	if err != nil {
		log.Printf("can not connect to storage: %s\n", err)
		return
	}
	servOrders := order.NewService(&storOrders, cover.NewService())
	servOrders.AddGivenOrdersGauge(givenOrdersCounter)
	servOrders.AddFailedRequestsCounter(failedOrderCounter)

	grpcMetrics := grpc_prometheus.NewServerMetrics()
	reg.MustRegister(grpcMetrics)

	serv := server.NewServer(service, producer, reg)
	serv.AddOrderService(servOrders)
	serv.AddGRPCMetrics(grpcMetrics)

	log.Println("Ready to run")

	if err := serv.RunGRPC(ctx, cfg); err != nil {
		log.Fatal(err)
	}
}
