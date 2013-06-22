package api

import (
	"time"
)

type HttpCallbackService struct {
	StartTime time.Time

	Callbacks *CallbackController
	Users     *UserController
}

func NewService(callbackCtlr *CallbackController, usersCtlr *UserController) *HttpCallbackService {
	return &HttpCallbackService{
		StartTime: time.Now(),
		Callbacks: callbackCtlr,
		Users:     usersCtlr,
	}
}
