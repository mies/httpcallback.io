package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/pjvds/httpcallback.io/api/host"
	"net/http"
	"os"
)

var (
	Address    = flag.String("address", "", "the address to host on")
	Port       = flag.Int("port", 8000, "the port to host on")
	ConfigPath = flag.String("config", "config.toml", "the path to the configuration file")
)

func main() {
	flag.Parse()
	InitLogging()

	Log.Info("Starting with config %v\n", *ConfigPath)
	config, err := host.OpenConfig(*ConfigPath)
	if err != nil {
		Log.Error("Unable to open config: %v", err.Error())
		return
	}

	apiHost := host.NewServer(config)

	//w := worker.NewCallbackWorker(100*time.Millisecond, repositoryFactory.CreateCallbackRepository())
	//w.Start()
	//Log.Notice("Started worker!")

	address := fmt.Sprintf("%v:%v", *Address, *Port)
	Log.Info("httpcallback.io now hosting at %v\n", address)
	if err := http.ListenAndServe(address, handlers.LoggingHandler(os.Stdout, apiHost)); err != nil {
		Log.Fatal(err)
	}
}
