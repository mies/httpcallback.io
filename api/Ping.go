package api

import (
	"net/http"
)

func Ping(response http.ResponseWriter, request *http.Request) {
	data := ResponseData{"message": "pong"}

	response.WriteHeader(http.StatusOK)
	response.Write(data.ToJson())
}
