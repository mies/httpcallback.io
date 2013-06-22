package data

import (
	"github.com/pjvds/httpcallback.io/model"
)

type UserRepository interface {
	List() ([]*model.User, error)
	Add(*model.User) error
}
