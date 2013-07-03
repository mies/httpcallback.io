package mvc

import (
	"encoding/json"
	"net/http"
	"testing"
)

type ResponseRecorder struct {
	header     http.Header
	StatusCode int
	Body       []byte
}

func NewResponseRecorder() *ResponseRecorder {
	return &ResponseRecorder{
		header:     make(http.Header, 0),
		StatusCode: http.StatusOK,
		Body:       make([]byte, 0),
	}
}

func (r *ResponseRecorder) Header() http.Header {
	return r.header
}

func (r *ResponseRecorder) WriteHeader(status int) {
	r.StatusCode = status
}

func (r *ResponseRecorder) Write(content []byte) (int, error) {
	for _, b := range content {
		r.Body = append(r.Body, b)
	}

	return len(content), nil
}

func TestJsonResultAddsContentTypeHeaderValue(t *testing.T) {
	result, _ := JsonResult(JsonDocument{
		"message": "foobar",
	})

	response := NewResponseRecorder()
	result.WriteResponse(response)

	contentType := response.Header().Get("Content-Type")
	if contentType == "" {
		t.Error("Content-Type header is not set")
	}
	if contentType != "application/json" {
		t.Errorf("Wrong Content-Type header set:\n\tactual: %v\n\rexpected: %v",
			contentType, "application/json")
	}
}

func TestJsonResultSerializesCorrectJson(t *testing.T) {
	result, _ := JsonResult(JsonDocument{
		"message": "foobar",
	})

	response := NewResponseRecorder()
	result.WriteResponse(response)

	var resultDoc JsonDocument
	if err := json.Unmarshal(response.Body, &resultDoc); err != nil {
		t.Fatalf("Invalid json in body: %v", err.Error())
	}
	if resultDoc["message"] != "foobar" {
		t.Fatalf("Wrong value received in result document:\n\tactual: %v\n\texpected: %v",
			resultDoc["message"], "foobar")
	}
}
