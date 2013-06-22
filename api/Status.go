package api

import (
	"net/http"
)

func Status(response api.ResponseWriter, request *http.Request) {
	switch request.Method {
	case "GET":
		response.Success(api.ResponseData{"message": "we are up & running!"})
	case "OPTIONS":
		response.SetHeader("Allow", "GET, OPTIONS")
	default:
		response.SetHeader("Allow", "GET, OPTIONS")
		response.ErrMethodNotAllowed()
	}
}
