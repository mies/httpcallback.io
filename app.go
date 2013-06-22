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
	DbUrl   = flag.String("db-url", "mongodb://foo:bar@dharma.mongohq.com:10039/httpcallback", "mongo url, including credentials, if needed.")
	DbName  = flag.String("db-name", "httpcallback", "mongo database name")
)

func main() {
	flag.Parse()
	mongoSession, err := mongo.Open(*DbUrl, *DbName)
	if err != nil {
		panic(err)
	}

	callbacksRepository := mongo.NewCallbackRepository(mongoSession)
	callbacksController := api.NewCallbackController(callbacksRepository)

	usersRepository := mongo.NewUserRepository(mongoSession)
	usersController := api.NewUserController(usersRepository)
	service := api.NewService(callbacksController, usersController)

	address := fmt.Sprintf("%s:%v", *Address, *Port)
	router := mux.NewRouter()

	siteRouter := router.Host("").Subrouter()
	siteRouter.Handle("/", http.FileServer(http.Dir("./site")))

	apiRouter := router.Host("api.{host:[a-z]+}").Subrouter()
	apiRouter.HandleFunc("/ping", HttpReponseWrapper(service.GetPing)).Methods("GET")
	apiRouter.HandleFunc("/users", HttpReponseWrapper(service.Users.ListUsers)).Methods("GET")
	apiRouter.HandleFunc("/callbacks", func(response http.ResponseWriter, req *http.Request) {
		result, err := service.Callbacks.ListCallbacks(req)
		WriteResultOrError(response, result, err)
	}).Methods("GET")
	apiRouter.HandleFunc("/callbacks", func(response http.ResponseWriter, req *http.Request) {
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
