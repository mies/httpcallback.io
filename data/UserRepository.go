package data

import (
	"github.com/pjvds/httpcallback.io/model"
)

type UserRepository interface {
	Add(*model.User) error
	Get(model.ObjectId) (*model.User, error)
	GetByUsernameAndPasswordHash(username string, passwordHash string) (*model.User, error)
	GetByUsernameAndAuthToken(username string, authToken model.AuthenticationToken) (*model.UserAuthInfo, error)
	List() ([]*model.User, error)
}
