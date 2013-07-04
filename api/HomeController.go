package api

import (
	"fmt"
	. "github.com/pjvds/httpcallback.io/mvc"
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

func (c *HomeController) HandleIndex(request *http.Request) ActionResult {
	up := time.Now().Sub(c.StartTime)
	uptime := fmt.Sprintf("%v:%v:%v", up.Hours(), up.Minutes(), up.Seconds())

	return JsonResult(&JsonDocument{
		"message": "welcome!",
		"uptime":  uptime,
	})
}

type PingResponse struct {
	Message string `json:"message"`
}

func (c *HomeController) HandlePing(req *http.Request) ActionResult {
	return JsonResult(&PingResponse{
		Message: "pong",
	})
}
