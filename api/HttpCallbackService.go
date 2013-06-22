package api

import (
	"time"
)

type HttpCallbackService struct {
	StartTime time.Time
}

func NewService() *HttpCallbackService {
	return &HttpCallbackService{
		StartTime: time.Now(),
	}
}
