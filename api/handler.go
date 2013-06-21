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
	SetHeader(string, string)
	Success(Response)
	ErrMethodNotAllowed()
}

func (writer HttpResponseWriter) SetHeader(key string, value string) {
	writer.Header().Set(key, value)
}

func (writer HttpResponseWriter) Success(response Response) {
	writer.Write(http.StatusOK, response)
}

func (writer HttpResponseWriter) ErrMethodNotAllowed() {
	writer.Write(http.StatusMethodNotAllowed, ResponseData{"message": "Method Not Allowed"})
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
