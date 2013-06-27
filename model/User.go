package model

import (
	"time"
)

type User struct {
	Id           ObjectId            `json:"id"`
	CreatedAt    time.Time           `json:"createAt"`
	Username     string              `json:"username"`
	PasswordHash string              `json:"passwordHash"`
	AuthToken    AuthenticationToken `json:"authToken"`
}
