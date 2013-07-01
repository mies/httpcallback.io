package api

import (
	"time"
)

type HttpCallbackService struct {
	StartTime time.Time
	Home      *HomeController
	Callbacks *CallbackController
	Users     *UserController
}

func NewService(callbackCtlr *CallbackController, usersCtlr *UserController) *HttpCallbackService {
	return &HttpCallbackService{
		StartTime: time.Now(),
		Home:      NewHomeController(),
		Callbacks: callbackCtlr,
		Users:     usersCtlr,
	}
}
