package api

import (
	"github.com/nu7hatch/gouuid"
	"net/http"
	"time"
)

type CallbackController struct {
}

func NewCallbackController() *CallbackController {
	return &CallbackController{}
}

type CallbackRequestArgs struct {
	When time.Time `json:"when"`
	Url  string    `json:"url"`
}

type CallbackRequestReply struct {
	Id string `json:"id"`
}

func (ctr *CallbackController) NewCallback(r *http.Request, args *CallbackRequestArgs) (*JsonResponse, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	reply := &CallbackRequestReply{
		Id: id.String(),
	}

	return JsonResult(reply)
}
