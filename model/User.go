package model

import (
	//"github.com/pjvds/httpcallback.io/security"
	"time"
)

type User struct {
	Id           ObjectId  `json:"id"`
	CreatedAt    time.Time `json:"createAt"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"passwordHash"`
	//AuthToken    security.AuthenticationToken `json:"authToken"`
}
