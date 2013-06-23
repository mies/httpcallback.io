package security

import (
	"crypto/sha256"
	"fmt"
	"github.com/nu7hatch/gouuid"
	"time"
)

func HashPassword(username string, password string, creationDate time.Time) string {
	data := []byte(fmt.Sprintf("%s:%s:%s", username, password, creationDate))

	hash := sha256.New()
	result := hash.Sum(data)
	return string(result)
}

type AuthenticationToken string

func NewAuthToken() AuthenticationToken {
	token, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}

	return AuthenticationToken(token.String())
}

func (token AuthenticationToken) String() string {
	return string(token)
}
