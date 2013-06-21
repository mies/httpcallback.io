package main

import (
	"fmt"
	"github.com/pjvds/httpcallback.io/api"
	"net/http"
)

func main() {
	address := ":8000"
	fmt.Println("Hosting at ", address)

	api.AddHandler("/status", Status)
	if err := http.ListenAndServe(address, nil); err != nil {
		fmt.Println("error: ", err)
	}
}

func Status(response api.ResponseWriter, request *http.Request) {
	response.SetHeader("Allow", "GET")

	if request.Method == "GET" {
		response.Success(api.ResponseData{"message": "we are up & running!"})
	} else {
		response.ErrMethodNotAllowed()
	}
}
