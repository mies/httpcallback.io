package api

import (
	"github.com/pjvds/httpcallback.io/data"
	"github.com/pjvds/httpcallback.io/model"
)

type Authenticator struct {
	userRepository data.UserRepository
}

func NewAuthenticator(userRepository data.UserRepository) *Authenticator {
	return &Authenticator{
		userRepository: userRepository,
	}
}

func (a Authenticator) Authenticate(username string, token string) (*model.UserAuthInfo, error) {
	return a.userRepository.GetByAuth(username, model.AuthenticationToken(token))
}

func (a Authenticator) Wrap(handler AuthenticatedRequestHandler) *AuthenticationHandler {
	return &AuthenticationHandler{
		authenticator: a,
		Handler:       handler,
	}
}
