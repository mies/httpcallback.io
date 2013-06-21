package api

import (
	"net/http"
)

type HttpResponseWriter struct {
	http.ResponseWriter
}

func (writer *HttpResponseWriter) Write(statusCode int, response Response) {
	data := response.ToJson()
	writer.Header().Set("Content-Type", "application/json")
	writer.ResponseWriter.Write(data)
	writer.WriteHeader(statusCode)
}

type RequestHandlerFunc func(ResponseWriter, *http.Request)

type ResponseWriter interface {
	Success(Response)
}

func (writer HttpResponseWriter) Success(response Response) {
	writer.Write(http.StatusOK, response)
}

func AddHandler(pattern string, handler RequestHandlerFunc) {
	http.HandleFunc(pattern, createHandler(handler))
}

func createHandler(handler RequestHandlerFunc) func(http.ResponseWriter, *http.Request) {
	return func(response http.ResponseWriter, request *http.Request) {
		wrapper := &HttpResponseWriter{response}
		handler(wrapper, request)
	}
}
