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

func (ctr *CallbackController) NewCallback(r *AuthenticatedRequest, args *model.CallbackRequest) (*JsonResponse, error) {
	id := model.NewObjectId()
	callback := model.NewCallback(id, r.UserId, args)
	err := ctr.callbacks.Add(callback)
	if err != nil {
		return nil, err
	}

	return JsonResult(&struct {
		Id model.ObjectId `bson:"_id,omitempty json:"id"`
	}{
		Id: callback.Id,
	})
}

func (ctr *CallbackController) ListCallbacks(r *AuthenticatedRequest) (*JsonResponse, error) {
	callbacks, err := ctr.callbacks.List(r.UserId)
	if err != nil {
		return nil, err
	} else {
		return JsonResult(callbacks)
	}
}
