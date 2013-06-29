package api

import (
	"github.com/pjvds/httpcallback.io/data"
	"github.com/pjvds/httpcallback.io/model"
	"labix.org/v2/mgo/bson"
	"net/http"
	"time"
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
	callback := model.Callback{
		Id:        bson.NewObjectId().String(),
		CreatedAt: time.Now(),
		Request:   args,
	}
	err := ctr.callbacks.Add(&callback)
	if err != nil {
		return nil, err
	}

	return JsonResult(&struct {
		Id string `json:"id"`
	}{
		Id: callback.Id,
	})
}

func (ctr *CallbackController) ListCallbacks(r *http.Request) (*JsonResponse, error) {
	callbacks, err := ctr.callbacks.List()
	if err != nil {
		return nil, err
	} else {
		return JsonResult(callbacks)
	}
}
