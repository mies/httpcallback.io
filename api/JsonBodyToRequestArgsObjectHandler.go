package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
)

type JsonBodyToRequestArgsObjectHandler struct {
	handlerType    reflect.Type
	argsObjectType reflect.Type
	handlerValue   reflect.Value
}

func validateHandler(handler interface{}) (bool, error) {
	handlerType := reflect.TypeOf(handler)

	// Must be a func
	if handlerType.Kind() != reflect.Func {
		return false, errors.New(fmt.Sprintf("invalid handler type %s, handler must be an func", handlerType.Kind().String()))
	}

	// Must have 2 in parameters
	if handlerType.NumIn() != 2 {
		return false, errors.New(fmt.Sprintf("handler does not have 2 in parameters, instead it has %v", handlerType.NumIn()))
	}
	// Must have 2 out parameters
	if handlerType.NumOut() != 2 {
		return false, errors.New(fmt.Sprintf("handler does not have 2 out parameters, instead it has %v", handlerType.NumOut()))
	}

	// First in parameter must be *http.Request

	// TODO: Add type check for both, auth and non auth request
	//requestType := reflect.TypeOf(&http.Request{})
	// if handlerType.In(0).AssignableTo(requestType) {
	// 	return false, errors.New(fmt.Sprintf("invalid argument type, first argument should be of kind %v, not %v",
	// 		requestType.String(), handlerType.In(0).String()))
	// }

	// Second in parameter must be *struct
	if handlerType.In(1).Kind() != reflect.Ptr || handlerType.In(1).Elem().Kind() != reflect.Struct {
		return false, errors.New(fmt.Sprintf("invalid argument type, second argument should be a pointer to an struct, not %v",
			handlerType.In(1).String()))
	}

	// First out parameter must be HttpResponse
	if !handlerType.Out(0).Implements(reflect.TypeOf((*HttpResponse)(nil)).Elem()) {
		return false, errors.New(fmt.Sprintf("invalid argument type, first out parameter of type %v should implement api.HttpResponse interface",
			handlerType.Out(0).String()))
	}

	// Second out parameter must be error
	if handlerType.Out(1) != reflect.TypeOf((*error)(nil)).Elem() {
		return false, errors.New(fmt.Sprintf("invalid argument type, second out parameter of type %v should implement error interface",
			handlerType.Out(1).String()))
	}

	return true, nil
}

func NewJsonBodyRequestArgsObjectHandler(handler interface{}) *JsonBodyToRequestArgsObjectHandler {
	if ok, err := validateHandler(handler); !ok {
		panic(err)
	}

	handlerType := reflect.TypeOf(handler)
	h := &JsonBodyToRequestArgsObjectHandler{
		handlerType:    handlerType,
		argsObjectType: handlerType.In(1),
		handlerValue:   reflect.ValueOf(handler),
	}

	Log.Debug("Created JsonBodyRequestAgsObjectHandler for func %v", h.handlerType.String())
	return h
}

func (h *JsonBodyToRequestArgsObjectHandler) invoke(request reflect.Value, args reflect.Value) (HttpResponse, error) {
	var result HttpResponse
	var err error

	results := h.handlerValue.Call([]reflect.Value{request, args})
	if !results[0].IsNil() {
		result = results[0].Interface().(HttpResponse)
	}
	if !results[1].IsNil() {
		err = results[1].Interface().(error)
	}

	return result, err
}

func (h *JsonBodyToRequestArgsObjectHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)

	argsObjectPtr := reflect.New(h.argsObjectType.Elem())
	argsObject := argsObjectPtr.Interface()

	if err := decoder.Decode(&argsObject); err != nil {
		Log.Warning("invalid body for request object type %v: %v", h.argsObjectType.Name(), err.Error())
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte(err.Error()))
		return
	}

	result, err := h.invoke(reflect.ValueOf(request), reflect.ValueOf(argsObject))

	if err != nil {
		Log.Error("error from handler %v: %v", h.handlerType.String(), err.Error())
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	result.WriteResponse(response)
}

func (h *JsonBodyToRequestArgsObjectHandler) ServeAuthHTTP(response http.ResponseWriter, request *AuthenticatedRequest) {
	decoder := json.NewDecoder(request.Body)

	argsObjectPtr := reflect.New(h.argsObjectType.Elem())
	argsObject := argsObjectPtr.Interface()

	if err := decoder.Decode(&argsObject); err != nil {
		Log.Warning("invalid body for request object type %v: %v", h.argsObjectType.Name(), err.Error())
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte(err.Error()))
		return
	}

	result, err := h.invoke(reflect.ValueOf(request), reflect.ValueOf(argsObject))

	if err != nil {
		Log.Error("error from handler %v: %v", h.handlerType.String(), err.Error())
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	result.WriteResponse(response)
}
