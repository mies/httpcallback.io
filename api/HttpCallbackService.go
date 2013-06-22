package api

import (
	"time"
)

type HttpCallbackService struct {
	StartTime time.Time

	Callbacks *CallbackController
}

func NewService(callbackCtlr *CallbackController) *HttpCallbackService {
	return &HttpCallbackService{
		StartTime: time.Now(),
		Callbacks: callbackCtlr,
	}
}
