package server

import (
	"context"
	"log"
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
	"gitlab.ozon.dev/mer_marat/homework/internal/model"
	"gitlab.ozon.dev/mer_marat/homework/internal/pkg/kafka"
)

type producer interface {
	SendSyncMessage(message *sarama.ProducerMessage) (partition int32, offset int64, err error)
}

type servicePickpoint interface {
	Read(ctx context.Context, id int64) (res *model.PickPoint, err error)
	Create(ctx context.Context, point *model.PickPoint) (res *model.PickPoint, err error)
	Update(ctx context.Context, point *model.PickPoint) (err error)
	Delete(ctx context.Context, id int64) (err error)
}

type server struct {
	servicePickpoint
	producer
	reg         *prometheus.Registry
	grpcMetrics *grpc_prometheus.ServerMetrics
}

func NewServer(servicePickpoint servicePickpoint, producer producer, reg *prometheus.Registry) server {
	return server{
		servicePickpoint: servicePickpoint,
		producer:         producer,
		reg:              reg,
	}
}

func (s *server) AddGRPCMetrics(grpcMetrics *grpc_prometheus.ServerMetrics) {
	s.grpcMetrics = grpcMetrics
}

func (s server) Run(ctx context.Context, cfg config.Config) error {
	sender := kafka.NewKafkaSender(s.producer, cfg.Kafka.Topic)
	handler := handler.NewHandler(s.servicePickpoint)
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
