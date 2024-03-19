package main

import (
	"context"
	"fmt"
	"homework2/pup/internal/pkg/db"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	queryParamKey = "point"
	port          = ":9000"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	database, err := db.NewDB(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	router := mux.NewRouter()
	router.HandleFunc("/pickpoint", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			fmt.Println("POST")
		case http.MethodPut:
			fmt.Println("PUT")
		default:
			fmt.Println("error")
		}
	})

	router.HandleFunc(fmt.Sprintf("/pickpoint/{%s:[A-z]+}", queryParamKey), func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			fmt.Println("GET")
		case http.MethodDelete:
			fmt.Println("DELETE")
		default:
			fmt.Println()
		}
	})

	http.Handle("/", router)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}
