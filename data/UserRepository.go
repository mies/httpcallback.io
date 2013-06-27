package data

import (
	"github.com/pjvds/httpcallback.io/model"
)

type UserRepository interface {
	Add(*model.User) error
	Get(model.ObjectId) (*model.User, error)
	GetByAuth(username string, authToken model.AuthenticationToken) (*model.User, error)
	List() ([]*model.User, error)
}
