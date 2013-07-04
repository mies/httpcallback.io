package api

import (
	. "github.com/pjvds/httpcallback.io/api/messages"
	"github.com/pjvds/httpcallback.io/data"
	"github.com/pjvds/httpcallback.io/model"
	. "github.com/pjvds/httpcallback.io/mvc"
)

type CallbackController struct {
	callbacks data.CallbackRepository
}

func NewCallbackController(callbacks data.CallbackRepository) *CallbackController {
	return &CallbackController{
		callbacks: callbacks,
	}
}

func (ctr *CallbackController) NewCallback(r *AuthenticatedRequest, args *NewCallbackRequest) ActionResult {
	callback := model.NewCallback(r.UserId, args.Url, args.When)

	if err := ctr.callbacks.Add(callback); err != nil {
		Log.Error("Unable to add callback to repository: %v", err.Error())
		return ErrorResult(err)
	}

	return JsonResult(&NewCallbackResponse{
		Id: callback.Id.String(),
	})
}

func (ctr *CallbackController) ListCallbacks(r *AuthenticatedRequest) ActionResult {
	callbacks, err := ctr.callbacks.List(r.UserId)

	if err != nil {
		Log.Error("Error while getting callback for user '%v': %v", r.UserId, err.Error())
		return ErrorResult(err)
	}

	return JsonResult(callbacks)
}
