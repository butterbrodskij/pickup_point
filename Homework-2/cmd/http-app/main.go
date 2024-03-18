package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	queryParamKey = "point"
	port          = ":9000"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/pickpoint", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			fallthrough
		case http.MethodPut:
			fallthrough
		default:
			fmt.Println()
		}
	})

	router.HandleFunc(fmt.Sprintf("/pickpoint/{%s:[A-z]+}", queryParamKey), func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			fallthrough
		case http.MethodDelete:
			fallthrough
		default:
			fmt.Println()
		}
	})

	http.Handle("/", router)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}
