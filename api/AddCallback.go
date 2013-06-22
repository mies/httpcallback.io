package api

import (
	"github.com/nu7hatch/gouuid"
	"net/http"
	"net/url"
	"time"
)

type CallbackRequestArgs struct {
	When time.Time `json:"when"`
	Url  url.URL   `json:"url"`
}

type CallbackRequestReply struct {
	Id *uuid.UUID `json:"id"`
}

func (s *HttpCallbackService) AddCalback(r *http.Request, args *CallbackRequestArgs, reply *CallbackRequestReply) (err error) {
	reply.Id, err = uuid.NewV4()
	return
}
