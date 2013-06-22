package api

import (
	"github.com/nu7hatch/gouuid"
	"net/http"
	"strconv"
	"time"
)

type CallbackController struct {
	callbacks []Callback
}

func NewCallbackController() *CallbackController {
	return &CallbackController{callbacks: make([]Callback, 0)}
}

type Callback struct {
	Id        string               `json:"id"`
	CreatedAt time.Time            `json:"createAt"`
	Request   *CallbackRequestArgs `json:"request"`
}

type CallbackRequestArgs struct {
	When time.Time `json:"when"`
	Url  string    `json:"url"`
}

type CallbackRequestReply struct {
	Id string `json:"id"`
}

func (ctr *CallbackController) ListCallbacks(r *http.Request) (*JsonResponse, error) {
	num := len(ctr.callbacks)
	print("callbacks: " + strconv.Itoa(num))
	return JsonResult(ctr.callbacks)
}

func (ctr *CallbackController) NewCallback(r *http.Request, args *CallbackRequestArgs) (*JsonResponse, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	callback := Callback{
		Id:        id.String(),
		CreatedAt: time.Now(),
		Request:   args,
	}
	ctr.callbacks = append(ctr.callbacks, callback)
	print("added!")
	print("callbacks: " + strconv.Itoa(len(ctr.callbacks)))

	reply := &CallbackRequestReply{
		Id: callback.Id,
	}

	return JsonResult(reply)
}
