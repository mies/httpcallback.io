package mvc

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

type RequestArgs struct{}

func TestNewingPanicsOnWrongKind(t *testing.T) {
	ok, err := validateHandler("iamastringnotahandler")

	if ok {
		t.Error("Handler of wrong kind should not be valid")
	}

	if err == nil {
		t.Fatal("Handler of wrong kind should return error")
	}

	expectedError := fmt.Sprint("invalid handler type string, handler must be an func")
	if !strings.Contains(err.Error(), expectedError) {
		t.Errorf("Unexpected error message: \n\tActual: %v\n\tExpected: %v", err.Error(), expectedError)
	}
}

func TestNewingPanicsOnWrongInParameterCount(t *testing.T) {
	ok, err := validateHandler(func() {})

	if ok {
		t.Error("Handler with wrong parameter count should not be valid")
	}

	if err == nil {
		t.Fatal("Handler with wring parameter count should return error")
	}

	expectedError := fmt.Sprint("handler does not have 2 in parameters, instead it has 0")
	if !strings.Contains(err.Error(), expectedError) {
		t.Errorf("Unexpected error message: \n\tActual: %v\n\tExpected: %v", err.Error(), expectedError)
	}
}

func TestNewingPanicsOnWrongOutParameterCount(t *testing.T) {
	ok, err := validateHandler(func(req *http.Request, args *RequestArgs) {})

	if ok {
		t.Error("Handler with wrong parameter count should not be valid")
	}

	if err == nil {
		t.Fatal("Handler with wring parameter count should return error")
	}

	expectedError := fmt.Sprint("handler does not have 1 out parameter, instead it has 0")
	if !strings.Contains(err.Error(), expectedError) {
		t.Errorf("Unexpected error message: \n\tActual: %v\n\tExpected: %v", err.Error(), expectedError)
	}
}

// TODO: Enable
func TestNewingPanicsOnWrongFirstInParameterType(t *testing.T) {
	t.Skipf("Temporary disble first parameter check")

	ok, err := validateHandler(func(s *string, args *RequestArgs) ActionResult {
		return nil
	})

	if ok {
		t.Error("Wrong hander should not be valid")
	}

	if err == nil {
		t.Fatal("Wrong hander should return error")
	}

	expectedError := fmt.Sprint("invalid argument type, first argument should be of kind *http.Request, not *string")
	if !strings.Contains(err.Error(), expectedError) {
		t.Errorf("Unexpected error message: \n\tActual: %v\n\tExpected: %v", err.Error(), expectedError)
	}
}

func TestNewingPanicsOnWrongSecondInParameterTypeNotPtr(t *testing.T) {
	ok, err := validateHandler(func(req *http.Request, args RequestArgs) ActionResult {
		return nil
	})

	if ok {
		t.Error("Wrong hander should not be valid")
	}

	if err == nil {
		t.Fatal("Wrong hander should return error")
	}

	expectedError := fmt.Sprint("invalid argument type, second argument should be a pointer to an struct, not mvc.RequestArgs")
	if !strings.Contains(err.Error(), expectedError) {
		t.Errorf("Unexpected error message: \n\tActual: %v\n\tExpected: %v", err.Error(), expectedError)
	}
}

func TestNewingPanicsOnWrongSecondInParameterTypeNotStruct(t *testing.T) {
	ok, err := validateHandler(func(req *http.Request, args *func()) ActionResult {
		return nil
	})

	if ok {
		t.Error("Wrong hander should not be valid")
	}

	if err == nil {
		t.Fatal("Wrong hander should return error")
	}

	expectedError := fmt.Sprint("invalid argument type, second argument should be a pointer to an struct, not *func()")
	if !strings.Contains(err.Error(), expectedError) {
		t.Errorf("Unexpected error message: \n\tActual: %v\n\tExpected: %v", err.Error(), expectedError)
	}
}

func TestNewingPanicsOnWrongFirstOutParameterTypeNotImplementingActionResult(t *testing.T) {
	ok, err := validateHandler(func(req *http.Request, args *RequestArgs) *int {
		return nil
	})

	if ok {
		t.Error("Wrong hander should not be valid")
	}

	if err == nil {
		t.Fatal("Wrong hander should return error")
	}

	expectedError := fmt.Sprint("invalid argument type, first out parameter of type *int should implement mvc.ActionResult interface")
	if !strings.Contains(err.Error(), expectedError) {
		t.Errorf("Unexpected error message: \n\tActual: %v\n\tExpected: %v", err.Error(), expectedError)
	}
}

func TestNewJsonBodyRequestArgsObjectHandlerSetType(t *testing.T) {
	h := NewJsonBodyRequestArgsObjectHandler(func(req *http.Request, args *RequestArgs) ActionResult {
		return nil
	})

	if h.argsObjectType != reflect.TypeOf(&RequestArgs{}) {
		t.Fatalf("unexpected argsObjectType %v", h.argsObjectType.String())
	}
}
