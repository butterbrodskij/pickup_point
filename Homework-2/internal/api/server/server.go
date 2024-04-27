package server

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"

	"github.com/IBM/sarama"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus"
	handler "gitlab.ozon.dev/mer_marat/homework/internal/api/handlers/pickpoint"
	"gitlab.ozon.dev/mer_marat/homework/internal/api/middleware"
	"gitlab.ozon.dev/mer_marat/homework/internal/api/router"
	"gitlab.ozon.dev/mer_marat/homework/internal/config"
	"gitlab.ozon.dev/mer_marat/homework/internal/metrics"
	"gitlab.ozon.dev/mer_marat/homework/internal/pkg/kafka"
	order_pb "gitlab.ozon.dev/mer_marat/homework/internal/pkg/pb/homework/orders/v1"
	pickpoint_pb "gitlab.ozon.dev/mer_marat/homework/internal/pkg/pb/homework/pickpoints/v1"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

type producer interface {
	SendSyncMessage(message *sarama.ProducerMessage) (partition int32, offset int64, err error)
}

type server struct {
	service_pickpoint pickpoint_pb.PickPointsServer
	service_order     order_pb.OrdersServer
	producer
	reg         *prometheus.Registry
	grpcMetrics *grpc_prometheus.ServerMetrics
}

func NewServer(service_pickpoint pickpoint_pb.PickPointsServer, producer producer, reg *prometheus.Registry) server {
	return server{
		service_pickpoint: service_pickpoint,
		producer:          producer,
		reg:               reg,
	}
}

func (s *server) AddOrderService(service_order order_pb.OrdersServer) {
	s.service_order = service_order
}

func (s *server) AddGRPCMetrics(grpcMetrics *grpc_prometheus.ServerMetrics) {
	s.grpcMetrics = grpcMetrics
}

func (s server) Run(ctx context.Context, cfg config.Config) error {
	sender := kafka.NewKafkaSender(s.producer, cfg.Kafka.Topic)
	handler := handler.NewHandler(s.service_pickpoint)
	authMiddleware := middleware.NewAuthMiddleware(cfg)
	logMiddleware := middleware.NewLogMiddleware(sender)

	router := router.MakeRouter(handler, authMiddleware, logMiddleware, cfg)
	errChan := make(chan error, 1)
	errSecChan := make(chan error, 1)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	wg := sync.WaitGroup{}
	wg.Add(2)
	defer wg.Wait()

	go func() {
		if err := metrics.Listen(":9095", s.reg); err != nil {
			log.Printf("metrics handling failed: %s", err)
		}
	}()

	secureServer := &http.Server{Addr: cfg.Server.SecurePort}
	server := &http.Server{Addr: cfg.Server.Port}
	http.Handle("/", router)
	go func() {
		defer wg.Done()
		if err := server.ListenAndServe(); err != nil {
			errChan <- err
		}
	}()
	go func() {
		defer wg.Done()
		if err := secureServer.ListenAndServeTLS(config.CertFile, config.KeyFile); err != nil {
			errSecChan <- err
		}
	}()

	select {
	case sig := <-sigChan:
		log.Printf("caught signal: %s\n", sig)
		err := server.Shutdown(ctx)
		errSec := secureServer.Shutdown(ctx)
		if err != nil {
			return err
		}
		if errSec != nil {
			return errSec
		}
	case err := <-errChan:
		secureServer.Shutdown(ctx)
		return err
	case err := <-errSecChan:
		server.Shutdown(ctx)
		return err
	}

	return nil
}

func (s server) RunGRPC(ctx context.Context, cfg config.Config) error {
	lis, err := net.Listen("tcp", ":9094")
	if err != nil {
		return err
	}
	defer lis.Close()

	grpcServer := grpc.NewServer()

	if s.grpcMetrics != nil {
		grpcServer = grpc.NewServer(
			grpc.StatsHandler(otelgrpc.NewServerHandler()),
			grpc.ChainUnaryInterceptor(
				s.grpcMetrics.UnaryServerInterceptor(),
			),
			grpc.ChainStreamInterceptor(
				s.grpcMetrics.StreamServerInterceptor(),
			),
		)
	}

	pickpoint_pb.RegisterPickPointsServer(grpcServer, s.service_pickpoint)
	if s.service_order != nil {
		order_pb.RegisterOrdersServer(grpcServer, s.service_order)
	}

	if s.grpcMetrics != nil {
		s.grpcMetrics.InitializeMetrics(grpcServer)
	}

	errChan := make(chan error, 1)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	wg := sync.WaitGroup{}
	wg.Add(1)
	defer wg.Wait()

	go func() {
		if err := metrics.Listen(":9095", s.reg); err != nil {
			log.Printf("metrics handling failed: %s", err)
		}
	}()

	go func() {
		defer wg.Done()
		if err := grpcServer.Serve(lis); err != nil {
			errChan <- err
		}
	}()

	select {
	case sig := <-sigChan:
		log.Printf("caught signal: %s\n", sig)
		err := lis.Close()
		if err != nil {
			return err
		}
	case err := <-errChan:
		return err
	}

	return nil
}
