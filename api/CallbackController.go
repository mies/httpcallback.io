package api

import (
	"github.com/pjvds/httpcallback.io/data"
	"github.com/pjvds/httpcallback.io/model"
)

type CallbackController struct {
	callbacks data.CallbackRepository
}

func NewCallbackController(callbacks data.CallbackRepository) *CallbackController {
	return &CallbackController{
		callbacks: callbacks,
	}
}

func (ctr *CallbackController) AddCallback(callback *model.Callback) {
	ctr.callbacks.Add(callback)
}
