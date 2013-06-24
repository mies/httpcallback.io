package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/pjvds/httpcallback.io/api"
	"github.com/pjvds/httpcallback.io/data"
	"github.com/pjvds/httpcallback.io/data/memory"
	"github.com/pjvds/httpcallback.io/data/mongo"
	"github.com/pjvds/httpcallback.io/model"
	"io/ioutil"
	"net/http"
)

var (
	Address    = flag.String("address", "", "the address to host on")
	Port       = flag.Int("port", 8000, "the port to host on")
	ConfigPath = flag.String("config", "config.toml", "the path to the configuration file")
)

func createRepositoryFactory(config *Configuration) (data.RepositoryFactory, error) {
	if config.Mongo.UseMongo {
		fmt.Println("Runnig with mongo data store")
		fmt.Printf("Connecting to mongo database %s... ", config.Mongo.DatabaseName)
		mongoSession, err := mongo.Open(config.Mongo.ServerUrl, config.Mongo.DatabaseName)
		if err != nil {
			fmt.Println("failed!")
			return nil, err
		}
		fmt.Println("succes!")
		return mongo.NewMgoRepositoryFactory(mongoSession), nil

	} else {
		fmt.Println("Runnig with inmemory data store")
		return memory.NewMemRepositoryFactory(), nil
	}
}

func main() {
	flag.Parse()
	fmt.Printf("Starting with config %s\n", *ConfigPath)
	config, err := OpenConfig(*ConfigPath)
	if err != nil {
		panic(err)
	}

	repositoryFactory, err := createRepositoryFactory(config)
	if err != nil {
		fmt.Println("[FATAL]" + err.Error())
		return
	}

	callbacksController := api.NewCallbackController(repositoryFactory.CreateCallbackRepository())
	usersController := api.NewUserController(repositoryFactory.CreateUserRepository())
	service := api.NewService(callbacksController, usersController)

	address := fmt.Sprintf("%s:%v", *Address, *Port)
	router := mux.NewRouter()

	siteRouter := router.Host(config.Host.Hostname).Subrouter()
	siteRouter.Handle("/", http.FileServer(http.Dir("./site")))

	apiRouter := router.Host("api." + config.Host.Hostname).Subrouter()
	apiRouter.HandleFunc("/ping", HttpReponseWrapper(service.GetPing)).Methods("GET")
	apiRouter.HandleFunc("/users", HttpReponseWrapper(service.Users.ListUsers)).Methods("GET")
	apiRouter.HandleFunc("/users", func(response http.ResponseWriter, req *http.Request) {
		decoder := json.NewDecoder(req.Body)
		var requestArgs api.AddUserRequest
		err = decoder.Decode(&requestArgs)
		if err != nil {
			fmt.Println("Error decoding body json to AddUserRequest: ", err)
			response.WriteHeader(http.StatusBadRequest)
			return
		}

		result, err := service.Users.AddUser(req, &requestArgs)
		WriteResultOrError(response, result, err)
	}).Methods("POST")
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
			fmt.Println("Error decoding body json to CallbackRequest: ", err)
			response.WriteHeader(http.StatusBadRequest)
			return
		}

		result, err := service.Callbacks.NewCallback(req, &args)
		WriteResultOrError(response, result, err)
	}).Methods("POST")

	fmt.Printf("httpcallback.io now hosting at %s\n", address)
	if err := http.ListenAndServe(address, router); err != nil {
		fmt.Println("fatal: ", err)
	}
}

func WriteResultOrError(w http.ResponseWriter, result api.HttpResponse, err error) {
	if err != nil {
		fmt.Println("Controller finished with error:", err)
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
