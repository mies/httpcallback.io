package api

import (
	"net/http"
	"time"
)

type HomeController struct {
	StartTime time.Time
}

func NewHomeController() *HomeController {
	return &HomeController{
		StartTime: time.Now(),
	}
}

func (c *HomeController) HandleIndex(request *http.Request) (ActionResult, error) {
	return JsonResult(&JsonDocument{
		"message": "welcome!",
		"uptime":  time.Now().Sub(c.StartTime).String(),
	})
}

type PingResponse struct {
	Message string `json:"message"`
}

func (c *HomeController) HandlePing(req *http.Request) (ActionResult, error) {
	return JsonResult(&PingResponse{
		Message: "pong",
	})
}
