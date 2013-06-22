package data

import (
	"github.com/pjvds/httpcallback.io/model"
)

type CallbackRepository interface {
	List() ([]*model.Callback, error)
	Add(*model.Callback) error
}
