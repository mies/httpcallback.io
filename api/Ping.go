package api

import (
	"net/http"
)

type PingResponse struct {
	Message string `json:"message"`
}

func (s *HttpCallbackService) GetPing(req *http.Request) (ActionResult, error) {
	return JsonResult(&PingResponse{
		Message: "pong",
	})
}
