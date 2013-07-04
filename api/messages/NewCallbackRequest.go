package messages

import (
	"time"
)

type NewCallbackRequest struct {
	When time.Time `json:"when"`
	Url  string    `json:"url"`
}
