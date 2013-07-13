package mvc

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"reflect"
)

type ParameterProvider interface {
	GetParamters(request *http.Request) map[string]string
}

type GorillaMuxVarsParameterProvider struct {
}

func (provider *GorillaMuxVarsParameterProvider) GetParamters(request *http.Request) map[string]string {
	return mux.Vars(request)
}

type JsonBodyToRequestArgsObjectHandler struct {
	handlerType       reflect.Type
	argsObjectType    reflect.Type
	handlerValue      reflect.Value
	parameterProvider *ParameterProvider
}

func validateHandler(handler interface{}) (bool, error) {
	handlerType := reflect.TypeOf(handler)

	// Must be a func
	if handlerType.Kind() != reflect.Func {
		return false, errors.New(fmt.Sprintf("invalid handler type %v, handler must be an func", handlerType.Kind().String()))
	}

	// Must have 2 in parameters
	if handlerType.NumIn() != 2 {
		return false, errors.New(fmt.Sprintf("handler does not have 2 in parameters, instead it has %v", handlerType.NumIn()))
	}

	// Must have 1 out parameters
	if handlerType.NumOut() != 1 {
		return false, errors.New(fmt.Sprintf("handler does not have 1 out parameter, instead it has %v", handlerType.NumOut()))
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

	// First out parameter must be ActionResult
	if !handlerType.Out(0).Implements(reflect.TypeOf((*ActionResult)(nil)).Elem()) {
		return false, errors.New(fmt.Sprintf("invalid argument type, first out parameter of type %v should implement mvc.ActionResult interface",
			handlerType.Out(0).String()))
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

func (h *JsonBodyToRequestArgsObjectHandler) invoke(request reflect.Value, args reflect.Value) ActionResult {
	results := h.handlerValue.Call([]reflect.Value{request, args})

	if results[0].IsNil() {
		panic("ActionResult of handler '" + h.handlerType.String() + "' is a nil reference")
	}

	return results[0].Interface().(ActionResult)
}

func (h *JsonBodyToRequestArgsObjectHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)

	argsObjectPtr := reflect.New(h.argsObjectType.Elem())
	argsObject := argsObjectPtr.Interface()

	var result ActionResult
	if err := decoder.Decode(&argsObject); err != nil {
		Log.Warning("invalid body for request object type %v: %v", h.argsObjectType.Name(), err.Error())

		result = ErrorResult(err)
	} else {
		result = h.invoke(reflect.ValueOf(request), reflect.ValueOf(argsObject))
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

	result := h.invoke(reflect.ValueOf(request), reflect.ValueOf(argsObject))
	result.WriteResponse(response)
}
