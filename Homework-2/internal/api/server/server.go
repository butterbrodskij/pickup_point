package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"

	"gitlab.ozon.dev/mer_marat/homework/internal/api/router"
	"gitlab.ozon.dev/mer_marat/homework/internal/config"
	"gitlab.ozon.dev/mer_marat/homework/internal/service/pickpoint"
)

type Server struct {
	service pickpoint.ServiceRepoInteface
}

func NewServer(service pickpoint.ServiceRepoInteface) Server {
	return Server{service: service}
}

func (s Server) Run(ctx context.Context, cfg config.Config, cancel context.CancelFunc) error {
	defer cancel()
	router := router.MakeRouter(ctx, s.service, cfg)
	errChan := make(chan error, 1)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	wg := sync.WaitGroup{}
	wg.Add(1)
	defer wg.Wait()

	server := &http.Server{Addr: cfg.Server.Port}
	http.Handle("/", router)
	go func() {
		defer wg.Done()
		if err := server.ListenAndServe(); err != nil {
			errChan <- err
		}
	}()

	select {
	case sig := <-sigChan:
		log.Printf("caught signal: %s\n", sig)
		if err := server.Shutdown(ctx); err != nil {
			return err
		}
	case err := <-errChan:
		return err
	}

	return nil
}
