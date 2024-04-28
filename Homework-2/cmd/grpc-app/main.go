package main

import (
	"context"
	"log"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus"
	grpc_order "gitlab.ozon.dev/mer_marat/homework/internal/api/grpc_handlers/orders/v1"
	grpc_pickpoint "gitlab.ozon.dev/mer_marat/homework/internal/api/grpc_handlers/pickpoints/v1"
	"gitlab.ozon.dev/mer_marat/homework/internal/api/server"
	"gitlab.ozon.dev/mer_marat/homework/internal/config"
	"gitlab.ozon.dev/mer_marat/homework/internal/metrics"
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

	pickpointCounter := metrics.PickpointCounter()
	requestPickpointMetrics := metrics.RequestPickpointHistogram()

	repo := postgres.NewRepo(database)
	cache := inmemorycache.NewInMemoryCache()
	defer cache.Close()
	//servicePoints := pickpoint.NewService(repo, redis, database)
	servicePoints := pickpoint.NewService(repo, cache, database)
	servicePoints.AddCounterMetric(&pickpointCounter)
	servicePoints.AddRequestHistogram(&requestPickpointMetrics)

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

	givenOrdersGauge := metrics.GivenOrdersGauge()
	failedOrderCounter := metrics.FailedOrderCounter()

	storOrders, err := storage.NewOrders("storage_orders.json")
	if err != nil {
		log.Printf("can not connect to storage: %s\n", err)
		return
	}
	serviceOrders := order.NewService(&storOrders, cover.NewService())
	serviceOrders.AddGivenOrdersGauge(&givenOrdersGauge)
	serviceOrders.AddFailedRequestsCounter(&failedOrderCounter)

	reg := prometheus.NewRegistry()

	grpcMetrics := grpc_prometheus.NewServerMetrics()
	reg.MustRegister(pickpointCounter, givenOrdersGauge, requestPickpointMetrics, failedOrderCounter, grpcMetrics)

	serv := server.NewServer(servicePoints, producer, reg)
	serv.AddGRPCMetrics(grpcMetrics)

	grpcServicePickpoint := grpc_pickpoint.NewGRPCPickpointService(servicePoints)
	grpcServiceOrder := grpc_order.NewGRPCOrderService(serviceOrders)

	log.Println("Ready to run")

	if err := serv.RunGRPC(ctx, cfg, grpcServicePickpoint, grpcServiceOrder); err != nil {
		log.Fatal(err)
	}
}
