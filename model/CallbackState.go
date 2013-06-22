package model

import (
	"time"
)

type Callback struct {
	Id        string           `json:"id"`
	CreatedAt time.Time        `json:"createAt"`
	Request   *CallbackRequest `json:"request"`
}

func NewCallback(id string, request *CallbackRequest) *Callback {
	return &Callback{
		Id:        id,
		CreatedAt: time.Now(),
		Request:   request,
	}
}
