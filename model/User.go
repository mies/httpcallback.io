package model

import (
	"time"
)

type User struct {
	Id           string    `json:"id"`
	CreatedAt    time.Time `json:"createAt"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"passwordHash"`
}
