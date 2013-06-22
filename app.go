package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/pjvds/httpcallback.io/api"
	"github.com/pjvds/httpcallback.io/data/mongo"
	"github.com/pjvds/httpcallback.io/model"
	"io/ioutil"
	"net/http"
)

var (
	Address = flag.String("address", "", "specifies the host address")
	Port    = flag.Int("port", 80, "specifies the host port")
	Mongo   = flag.String("mongo", "mongodb://website:Alien4LifeAndBeyond@dharma.mongohq.com:10039/httpcallback", "specifies mongo url")
)

func main() {
	flag.Parse()
	callbacksRepository, err := mongo.NewCallbackRepository(*Mongo)
	if err != nil {
		panic(err)
	}
	callbacksController := api.NewCallbackController(callbacksRepository)
	service := api.NewService(callbacksController)

	address := fmt.Sprintf("%s:%v", *Address, *Port)
	router := mux.NewRouter()
	router.Headers("Content-Type", "application/json")
	router.HandleFunc("/ping", HttpReponseWrapper(service.GetPing)).Methods("GET")
	router.HandleFunc("/callbacks", func(response http.ResponseWriter, req *http.Request) {
		result, err := service.Callbacks.ListCallbacks(req)
		WriteResultOrError(response, result, err)
	}).Methods("GET")
	router.HandleFunc("/callbacks", func(response http.ResponseWriter, req *http.Request) {
		data, err := ioutil.ReadAll(req.Body)
		if err != nil {
			panic(err)
		}

		var args model.CallbackRequest
		err = json.Unmarshal(data, &args)
		if err != nil {
			fmt.Println("ERROR:", err)
			response.WriteHeader(http.StatusBadRequest)
			return
		}

		result, err := service.Callbacks.NewCallback(req, &args)
		WriteResultOrError(response, result, err)
	}).Methods("POST")

	fmt.Println("Hosting at ", address)
	if err := http.ListenAndServe(address, router); err != nil {
		fmt.Println("fatal: ", err)
	}
}

func WriteResultOrError(w http.ResponseWriter, result api.HttpResponse, err error) {
	if err != nil {
		fmt.Println("ERROR:", err)
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		result.WriteResponse(w)
	}
}

func HttpReponseWrapper(handler func(*http.Request) (*api.JsonResponse, error)) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		result, err := handler(req)
		WriteResultOrError(res, result, err)
	}
}
