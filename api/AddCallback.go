package api

import (
	"github.com/nu7hatch/gouuid"
	. "github.com/pjvds/httpcallback.io/model"
	"net/http"
	"time"
)

type NewCallbackReply struct {
	Id string `json:"id"`
}

func (ctr *CallbackController) NewCallback(r *http.Request, args *CallbackRequest) (*JsonResponse, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	callback := Callback{
		Id:        id.String(),
		CreatedAt: time.Now(),
		Request:   args,
	}
	ctr.AddCallback(&callback)

	reply := &NewCallbackReply{
		Id: callback.Id,
	}

	return JsonResult(reply)
}

func (ctr *CallbackController) ListCallbacks(r *http.Request) (*JsonResponse, error) {
	return JsonResult(ctr.callbacks.List())
}
