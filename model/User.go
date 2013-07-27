package model

import (
	"github.com/pjvds/go-cqrs/sourcing"
	"time"
)

type UserCreated struct {
	Id           ObjectId  `json:"id"`
	Email        string    `json:"email"`
	CreatedAt    time.Time `json:"createAt"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"passwordHash"`
}

type UserAuthTokenAdded struct {
	AuthToken AuthenticationToken `json:"autoToken"`
}

type UserAuthInfo struct {
	UserId   ObjectId `json:"userId"`
	Username string   `json:"username"`
}

type User struct {
	sourcer sourcing.EventSource

	Id           ObjectId            `bson:"_id,omitempty" json:"id"`
	Email        string              `json:"email"`
	CreatedAt    time.Time           `json:"createAt"`
	Username     string              `json:"username"`
	PasswordHash string              `json:"passwordHash"`
	AuthToken    AuthenticationToken `json:"authToken"`
}

func NewUser(id ObjectId, createdAt time.Time, username, passwordHash, email string) *User {
	user := new(User)
	user.sourcer = sourcing.AttachNew(user)

	user.sourcer.Apply(UserCreated{
		Id:           id,
		Email:        email,
		Username:     username,
		PasswordHash: passwordHash,
		CreatedAt:    createdAt,
	})

	return user
}

func (u *User) AddAuthToken(token AuthenticationToken) {
	u.sourcer.Apply(UserAuthTokenAdded{
		AuthToken: token,
	})
}

func (u *User) HandleUserCreated(e UserCreated) {
	u.Id = e.Id
	u.Email = e.Email
	u.CreatedAt = e.CreatedAt
	u.Username = e.Username
	u.PasswordHash = e.PasswordHash
}

func (u *User) HandleUserAuthTokenAdded(e UserAuthTokenAdded) {
	u.AuthToken = e.AuthToken
}
