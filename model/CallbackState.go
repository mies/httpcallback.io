package model

import (
	"time"
)

type Callback struct {
	Id           ObjectId           `json:"id"`
	UserId       ObjectId           `json:"userId"`
	CreatedAt    time.Time          `json:"createAt"`
	Request      *CallbackRequest   `json:"request"`
	Attempts     []*CallbackAttempt `json:"attempts"`
	AttemptCount int                `json:"attemptCount"`
}

type CallbackEntry struct {
	Id   ObjectId  `json:"id"`
	When time.Time `json:"when"`
}

type CallbackAttempt struct {
	Id        ObjectId
	Timestamp time.Time
	Success   bool
	Message   string
	Response  *HttpResponseInfo
}

type HttpResponseInfo struct {
	HttpStatusCode int
	HttpStatusText string
	ResponseTime   time.Duration
}

func NewCallback(id ObjectId, userId ObjectId, request *CallbackRequest) *Callback {
	return &Callback{
		Id:        id,
		UserId:    userId,
		CreatedAt: time.Now(),
		Request:   request,
	}
}
