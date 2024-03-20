package main

import (
	"context"
	"homework2/pup/internal/api/router"
	"homework2/pup/internal/api/server"
	"homework2/pup/internal/pkg/db"
	"homework2/pup/internal/pkg/repository/postgres"
	"log"
	"net/http"
)

const (
	port = ":9000"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	database, err := db.NewDB(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	repo := postgres.NewRepo(database)
	serv := server.NewServer(repo)
	router := router.MakeRouter(ctx, serv)

	http.Handle("/", router)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}
