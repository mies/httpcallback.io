package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/pjvds/httpcallback.io/api"
	"net/http"
)

func main() {
	address := ":8000"
	router := mux.NewRouter()
	router.Headers("Content-Type", "application/json")
	router.HandleFunc("/ping", api.Ping).Methods("GET")

	fmt.Println("Hosting at ", address)
	if err := http.ListenAndServe(address, router); err != nil {
		fmt.Println("fatal: ", err)
	}
}
