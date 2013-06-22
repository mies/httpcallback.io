package api

import (
	"github.com/nu7hatch/gouuid"
	"net/http"
	"net/url"
	"time"
)

type CallbackController struct {
}

type CallbackRequestArgs struct {
	When time.Time `json:"when"`
	Url  url.URL   `json:"url"`
}

type CallbackRequestReply struct {
	Id *uuid.UUID `json:"id"`
}

func (ctr *CallbackController) AddCalback(r *http.Request, args *CallbackRequestArgs) (*JsonResponse, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	reply := &CallbackRequestReply{
		Id: id,
	}

	return JsonResult(reply)
}
