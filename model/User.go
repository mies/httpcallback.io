package model

import (
	"time"
)

type UserAuthInfo struct {
	UserId   ObjectId `json:"userId"`
	Username string   `json:"username"`
}

type User struct {
	Id           ObjectId            `bson:"_id,omitempty json:"id"`
	CreatedAt    time.Time           `json:"createAt"`
	Username     string              `json:"username"`
	PasswordHash string              `json:"passwordHash"`
	AuthToken    AuthenticationToken `json:"authToken"`
}
