package data

import (
	"github.com/pjvds/httpcallback.io/model"
)

type UserRepository interface {
	Get(model.ObjectId) (*model.User, error)
	List() ([]*model.User, error)
	Add(*model.User) error
}
