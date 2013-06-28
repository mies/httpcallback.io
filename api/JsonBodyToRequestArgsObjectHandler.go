package api

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
)

type JsonBodyToRequestArgsObjectHandler struct {
	handlerType    reflect.Type
	argsObjectType reflect.Type
}

func validateHandler(handler interface{}) (bool, error) {
	handlerType := reflect.TypeOf(handler)
	if handlerType.Kind() != reflect.Func {
		return false, errors.New(fmt.Sprintf("invalid handler type %s, handler must be an func", handlerType.Kind().String()))
	}
	if handlerType.NumIn() != 2 {
		return false, errors.New(fmt.Sprintf("handler does not have 2 in parameters, instead it has %v", handlerType.NumIn()))
	}

	expectedFirstArgType := reflect.TypeOf(http.Request{})
	if handlerType.In(0) != expectedFirstArgType {
		return false, errors.New(fmt.Sprintf("invalid argument type, first argument should be of type %v, not %v",
			expectedFirstArgType.Name(), handlerType.In(0).Name()))
	}

	if handlerType.In(1).Kind() != reflect.Ptr {
		return false, errors.New(fmt.Sprintf("invalid argument type, second argument should be of kind ptr, not %v",
			handlerType.In(1).Kind().String()))
	}

	return true, nil
}

func NewJsonBodyRequestArgsObjectHandler(handler interface{}) *JsonBodyToRequestArgsObjectHandler {
	if ok, err := validateHandler(handler); !ok {
		panic(err)
	}

	handlerType := reflect.TypeOf(handler)
	return &JsonBodyToRequestArgsObjectHandler{
		handlerType:    handlerType,
		argsObjectType: handlerType.In(1),
	}
}

// func (h *JsonBodyToRequestArgsObjectHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
// 	decoder := json.NewDecoder(request.Body)

// 	argsObjectPtr := reflect.New(h.argsObjectType)
// 	argsObject := argsObjectPtr.Elem().Interface()

// 	err := decoder.Decode(&argsObjectPtr)
// 	Log.Warning("invalid body for request object type %v: %v", h.argsObjectType.Name(), err.Error())

// }
