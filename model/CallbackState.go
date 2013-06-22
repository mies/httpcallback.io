package model

import (
	"time"
)

type Callback struct {
	Id        string           `json:"id"`
	CreatedAt time.Time        `json:"createAt"`
	Request   *CallbackRequest `json:"request"`
}
