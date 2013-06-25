package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type PingResponse struct {
	Message string `json:"message"`
}

type HttpResponse interface {
	WriteResponse(http.ResponseWriter)
}

func (s *HttpCallbackService) GetPing(req *http.Request) (*JsonResponse, error) {
	return JsonResult(&PingResponse{
		Message: "pong",
	})
}

func JsonResult(result interface{}) (*JsonResponse, error) {
	data, err := json.Marshal(result)
	if err != nil {
		fmt.Printf("Unable to marshal object (%+v) to json: %s", result, err.Error())
		return nil, err
	} else {
		return &JsonResponse{
			Data: data,
		}, nil
	}
}

type JsonResponse struct {
	Data []byte
}

func (j *JsonResponse) WriteResponse(response http.ResponseWriter) {
	response.Header().Set("Content-Type", "application/json")
	response.Write(j.Data)
}
