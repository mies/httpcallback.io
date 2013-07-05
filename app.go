package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/pjvds/httpcallback.io/api/host"
	"github.com/pjvds/httpcallback.io/data"
	"github.com/pjvds/httpcallback.io/data/memory"
	"github.com/pjvds/httpcallback.io/data/mongo"
	"github.com/pjvds/httpcallback.io/worker"
	"net/http"
	"os"
	"time"
)

var (
	Address    = flag.String("address", "", "the address to host on")
	Port       = flag.Int("port", 8000, "the port to host on")
	ConfigPath = flag.String("config", "config.toml", "the path to the configuration file")
)

func createRepositoryFactory(config *Configuration) (data.RepositoryFactory, error) {
	if config.Mongo.UseMongo {
		Log.Debug("Running with mongo data store")
		Log.Debug("Connecting to mongo database %v", config.Mongo.DatabaseName)
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
	InitLogging()

	Log.Info("Starting with config %v\n", *ConfigPath)
	config, err := OpenConfig(*ConfigPath)
	if err != nil {
		Log.Error("Unable to open config: %v", err.Error())
		return
	}

	repositoryFactory, err := createRepositoryFactory(config)
	if err != nil {
		Log.Fatal("[FATAL] Could not create repository factory: " + err.Error())
	}

	apiHost := host.NewServer(repositoryFactory)

	w := worker.NewCallbackWorker(100*time.Millisecond, repositoryFactory.CreateCallbackRepository())
	w.Start()
	Log.Notice("Started worker!")

	address := fmt.Sprintf("%v:%v", *Address, *Port)
	Log.Info("httpcallback.io now hosting at %v\n", address)
	if err := http.ListenAndServe(address, handlers.LoggingHandler(os.Stdout, apiHost)); err != nil {
		Log.Fatal(err)
	}
}
