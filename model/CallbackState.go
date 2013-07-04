package model

import (
	"time"
)

type Callback struct {
	Id                   ObjectId           `bson:"_id,omitempty" json:"id"`
	UserId               ObjectId           `json:"userId"`
	CreatedAt            time.Time          `json:"createAt"`
	Request              *CallbackRequest   `json:"request"`
	Attempts             []*CallbackAttempt `json:"attempts"`
	AttemptCount         int                `json:"attemptCount"`
	NextAttemptTimeStamp time.Time          `json:"nextAttemptTimestamp"`
	Finished             bool               `json:"finished"`
}

func NewCallback(userId ObjectId, url string, when time.Time) *Callback {
	return &Callback{
		Id:        NewObjectId(),
		UserId:    userId,
		CreatedAt: time.Now(),
		Request: &CallbackRequest{
			Url:  url,
			When: when,
		},
	}
}

type CallbackRequest struct {
	Url  string    `json:"url"`
	When time.Time `json:"when"`
}

type CallbackAttempt struct {
	Id        ObjectId
	Timestamp time.Time
	Success   bool
	Response  *CallbackHttpResponseInfo
}

type CallbackHttpResponseInfo struct {
	HttpStatusCode int
	HttpStatusText string
	ResponseTime   time.Duration
}
