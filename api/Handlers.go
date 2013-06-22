package api

import (
	"net/http"
)

type RequestHandler func(ResponseWriter, *http.Request)

func AddHandler(pattern string, handler RequestHandler) {
	http.HandleFunc(pattern, createHandler(handler))
}

func createHandler(handler RequestHandler) func(http.ResponseWriter, *http.Request) {
	return func(response http.ResponseWriter, request *http.Request) {
		wrapper := &HttpResponseWriter{response}
		handler(wrapper, request)
	}
}
