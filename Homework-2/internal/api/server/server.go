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
	"github.com/prometheus/client_golang/prometheus"
	"gitlab.ozon.dev/mer_marat/homework/internal/api/middleware"
	"gitlab.ozon.dev/mer_marat/homework/internal/api/router"
	"gitlab.ozon.dev/mer_marat/homework/internal/config"
	"gitlab.ozon.dev/mer_marat/homework/internal/metrics"
	"gitlab.ozon.dev/mer_marat/homework/internal/pkg/kafka"
	pickpoint_pb "gitlab.ozon.dev/mer_marat/homework/internal/pkg/pb"
	"gitlab.ozon.dev/mer_marat/homework/tests/dummy"
	"google.golang.org/grpc"
)

/*
type service interface {
	Create(context.Context, *model.PickPoint) (*model.PickPoint, error)
	Read(context.Context, int64) (*model.PickPoint, error)
	Update(context.Context, *model.PickPoint) error
	Delete(context.Context, int64) error
}
*/

type producer interface {
	SendSyncMessage(message *sarama.ProducerMessage) (partition int32, offset int64, err error)
}

type server struct {
	service pickpoint_pb.PickPointsServer
	producer
	reg prometheus.Gatherer
}

func NewServer(service pickpoint_pb.PickPointsServer, producer producer, reg prometheus.Gatherer) server {
	return server{
		service:  service,
		producer: producer,
		reg:      reg,
	}
}

func (s server) Run(ctx context.Context, cfg config.Config) error {
	sender := kafka.NewKafkaSender(s.producer, cfg.Kafka.Topic)
	//handler := handler.NewHandler(s.service)
	handler := dummy.NewHandlerApi()
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
	pickpoint_pb.RegisterPickPointsServer(grpcServer, s.service)

	go func() {
		if err := metrics.Listen(":9095", s.reg); err != nil {
			log.Printf("metrics handling failed: %s", err)
		}
	}()

	return grpcServer.Serve(lis)
}
