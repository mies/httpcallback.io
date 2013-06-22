package model

import (
	"time"
)

type CallbackRequest struct {
	When time.Time `json:"when"`
	Url  string    `json:"url"`
}
