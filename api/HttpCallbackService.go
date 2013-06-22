package api

import (
	"time"
)

type HttpCallbackService struct {
	StartTime time.Time

	Callbacks *CallbackController
}

func NewService() *HttpCallbackService {
	return &HttpCallbackService{
		StartTime: time.Now(),
		Callbacks: NewCallbackController(),
	}
}
