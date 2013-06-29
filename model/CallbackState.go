package model

import (
	"time"
)

type Callback struct {
	Id        ObjectId         `json:"id"`
	UserId    ObjectId         `json:"userId"`
	CreatedAt time.Time        `json:"createAt"`
	Request   *CallbackRequest `json:"request"`
}

func NewCallback(id ObjectId, userId ObjectId, request *CallbackRequest) *Callback {
	return &Callback{
		Id:        id,
		UserId:    userId,
		CreatedAt: time.Now(),
		Request:   request,
	}
}
