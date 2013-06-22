package api

import (
	"github.com/pjvds/httpcallback.io/model"
)

type CallbackController struct {
	callbacks []*model.Callback
}

func NewCallbackController() *CallbackController {
	return &CallbackController{}
}

func (ctr *CallbackController) AddCallback(callback *model.Callback) {
	ctr.callbacks = append(ctr.callbacks, callback)
}
