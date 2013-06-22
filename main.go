package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/pjvds/httpcallback.io/api"
	"net/http"
)

func main() {
	service := api.NewService()

	address := ":8000"
	router := mux.NewRouter()
	router.Headers("Content-Type", "application/json")
	router.HandleFunc("/ping", HttpReponseWrapper(service.GetPing)).Methods("GET")

	fmt.Println("Hosting at ", address)
	if err := http.ListenAndServe(address, router); err != nil {
		fmt.Println("fatal: ", err)
	}
}

func HttpReponseWrapper(handler func(*http.Request) (*api.JsonResponse, error)) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		w, err := handler(req)
		if err != nil {
			fmt.Println("ERROR:", err)
			res.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteResponse(res)
		}
	}
}
