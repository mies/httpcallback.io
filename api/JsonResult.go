package api

import (
	"encoding/json"
	"net/http"
)

type JsonResponse struct {
	Data []byte
}

func JsonResult(result interface{}) (*JsonResponse, error) {
	data, err := json.Marshal(result)
	if err != nil {
		Log.Error("Unable to marshal object (%+v) to json: %s", result, err.Error())
		return nil, err
	} else {
		return &JsonResponse{
			Data: data,
		}, nil
	}
}

func (j *JsonResponse) WriteResponse(response http.ResponseWriter) {
	response.Header().Set("Content-Type", "application/json")
	response.Write(j.Data)
}
