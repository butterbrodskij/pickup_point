package main

import (
	"context"
	"log"

	"gitlab.ozon.dev/mer_marat/homework/internal/api/server"
	"gitlab.ozon.dev/mer_marat/homework/internal/config"
	"gitlab.ozon.dev/mer_marat/homework/internal/pkg/db"
	"gitlab.ozon.dev/mer_marat/homework/internal/storage/postgres"
)

func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	database, err := db.NewDB(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	repo := postgres.NewRepo(database)
	serv := server.NewServer(repo)

	if err := serv.Run(ctx, cfg, cancel); err != nil {
		log.Fatal(err)
	}
}
