package main

import (
	"context"
	"log"
	"net/http"

	"gitlab.ozon.dev/mer_marat/homework/internal/api/router"
	"gitlab.ozon.dev/mer_marat/homework/internal/api/server"
	"gitlab.ozon.dev/mer_marat/homework/internal/pkg/db"
	"gitlab.ozon.dev/mer_marat/homework/internal/pkg/repository/postgres"
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
