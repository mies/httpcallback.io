// +build go1.1
package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/pjvds/httpcallback.io/api"
	"github.com/pjvds/httpcallback.io/data"
	"github.com/pjvds/httpcallback.io/data/memory"
	"github.com/pjvds/httpcallback.io/data/mongo"
	"net/http"
)

var (
	Address    = flag.String("address", "", "the address to host on")
	Port       = flag.Int("port", 8000, "the port to host on")
	ConfigPath = flag.String("config", "config.toml", "the path to the configuration file")
)

func createRepositoryFactory(config *Configuration) (data.RepositoryFactory, error) {
	if config.Mongo.UseMongo {
		Log.Debug("Running with mongo data store")
		Log.Debug("Connecting to mongo database %s", config.Mongo.DatabaseName)
		mongoSession, err := mongo.Open(config.Mongo.ServerUrl, config.Mongo.DatabaseName)
		if err != nil {
			Log.Error("Unable to connect to mongo:", err)
			return nil, err
		}
		Log.Debug("Connected succesfully")
		return mongo.NewMgoRepositoryFactory(mongoSession), nil

	} else {
		Log.Debug("Runnig with inmemory data store")
		return memory.NewMemRepositoryFactory(), nil
	}
}

func main() {
	flag.Parse()
	Log.Info("Starting with config %s\n", *ConfigPath)
	config, err := OpenConfig(*ConfigPath)
	if err != nil {
		panic(err)
	}

	repositoryFactory, err := createRepositoryFactory(config)
	if err != nil {
		Log.Fatal("[FATAL] Could not create repository factory: " + err.Error())
	}

	callbacksController := api.NewCallbackController(repositoryFactory.CreateCallbackRepository())
	usersController := api.NewUserController(repositoryFactory.CreateUserRepository())
	service := api.NewService(callbacksController, usersController)

	address := fmt.Sprintf("%s:%v", *Address, *Port)
	router := mux.NewRouter()

	siteRouter := router.Host(config.Host.Hostname).Subrouter()
	siteRouter.Handle("/", http.FileServer(http.Dir("./site")))

	apiRouter := router.Host("api." + config.Host.Hostname).Subrouter()
	apiPostRouter := apiRouter.Methods("POST").Subrouter()
	apiGetRouter := apiRouter.Methods("GET").Subrouter()

	apiGetRouter.HandleFunc("/ping", HttpReponseWrapper(service.GetPing))
	apiGetRouter.HandleFunc("/user/{id}", func(response http.ResponseWriter, req *http.Request) {
		Log.Info("[%v] %v\n", req.Method, req.URL)
		var result api.HttpResponse
		var err error

		userId, ok := mux.Vars(req)["id"]
		if !ok {
			Log.Warning("id parameter not given, return 404 not found.")
			result = api.NewHttpStatusCodeResult(http.StatusNotFound)
		} else {
			requestArgs := &api.GetUserRequestArgs{
				UserId: userId,
			}

			Log.Debug("Handing request to GetUser with %+v", requestArgs)
			result, err = service.Users.GetUser(req, requestArgs)
		}
		WriteResultOrError(response, result, err)
	})

	addUserHandler := api.NewJsonBodyRequestArgsObjectHandler(service.Users.AddUser)
	apiPostRouter.HandleFunc("/users", func(response http.ResponseWriter, req *http.Request) {
		Log.Info("[%v] %v\n", req.Method, req.URL)
		addUserHandler.ServeHTTP(response, req)
	})

	apiGetRouter.HandleFunc("/callbacks", func(response http.ResponseWriter, req *http.Request) {
		result, err := service.Callbacks.ListCallbacks(req)
		WriteResultOrError(response, result, err)
	})

	addCallbackHandler := api.NewJsonBodyRequestArgsObjectHandler(service.Callbacks.NewCallback)
	apiPostRouter.HandleFunc("/callbacks", func(response http.ResponseWriter, req *http.Request) {
		Log.Info("[%v] %v\n", req.Method, req.URL)
		addCallbackHandler.ServeHTTP(response, req)
	})

	Log.Info("httpcallback.io now hosting at %s\n", address)
	if err := http.ListenAndServe(address, router); err != nil {
		Log.Fatal(err)
	}
}

func WriteResultOrError(w http.ResponseWriter, result api.HttpResponse, err error) {
	if err != nil {
		Log.Debug("Controller finished with error: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		result.WriteResponse(w)
	}
}

func HttpReponseWrapper(handler func(*http.Request) (api.HttpResponse, error)) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		fmt.Printf("[%v] %v\n", req.Method, req.URL)
		result, err := handler(req)
		req.Body.Close()

		WriteResultOrError(res, result, err)
	}
}
