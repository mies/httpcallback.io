package data

import (
	"github.com/pjvds/httpcallback.io/model"
)

type CallbackRepository interface {
	List() []*model.Callback
	Add(*model.Callback)
}
