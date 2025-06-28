package server

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"

	"gitlab.ozon.dev/mer_marat/homework/internal/config"
	"gitlab.ozon.dev/mer_marat/homework/internal/metrics"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

type appHandler interface {
	RegisterGRPC(server grpc.ServiceRegistrar)
}

func (s server) RunGRPC(ctx context.Context, cfg config.Config, appHandlers ...appHandler) error {
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

	for _, handler := range appHandlers {
		handler.RegisterGRPC(grpcServer)
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
