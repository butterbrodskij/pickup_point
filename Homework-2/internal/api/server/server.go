package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"

	handler "gitlab.ozon.dev/mer_marat/homework/internal/api/handlers/pickpoint"
	"gitlab.ozon.dev/mer_marat/homework/internal/api/router"
	"gitlab.ozon.dev/mer_marat/homework/internal/config"
	"gitlab.ozon.dev/mer_marat/homework/internal/model"
)

type service interface {
	Create(context.Context, *model.PickPoint) (*model.PickPoint, error)
	Read(context.Context, int64) (*model.PickPoint, error)
	Update(context.Context, *model.PickPoint) error
	Delete(context.Context, int64) error
}

type server struct {
	service service
}

func NewServer(service service) server {
	return server{service: service}
}

func (s server) Run(ctx context.Context, cfg config.Config) error {
	router := router.MakeRouter(handler.NewHandler(s.service), cfg)
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
