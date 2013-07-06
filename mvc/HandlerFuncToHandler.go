package mvc

import (
	"net/http"
)

type HandlerFuncHandler http.HandlerFunc

func HandlerFuncToHandler(handlerFunc http.HandlerFunc) HandlerFuncHandler {
	return HandlerFuncHandler(handlerFunc)
}

func (h HandlerFuncHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	h(response, request)
}
