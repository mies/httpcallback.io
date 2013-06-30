package api

import (
	"net/http"
)

type HttpStatusCodeResult struct {
	StatusCode int
}

func NewHttpStatusCodeResult(statusCode int) *HttpStatusCodeResult {
	return &HttpStatusCodeResult{
		StatusCode: statusCode,
	}
}

func (h *HttpStatusCodeResult) WriteResponse(response http.ResponseWriter) {
	response.WriteHeader(h.StatusCode)
}
