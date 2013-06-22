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
	writer.WriteHeader(statusCode)
	writer.ResponseWriter.Write(data)
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
