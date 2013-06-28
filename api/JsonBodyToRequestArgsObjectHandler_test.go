package api

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
)

func HandlerWithWrongInParameterCount() (HttpResponse, error) {
	return nil, nil
}

type RequestArgs struct{}

func HandlerWithWrongOutParameterCount(req http.Request, args *RequestArgs) {
}

func TestNewingPanicsOnWrongKind(t *testing.T) {
	ok, err := validateHandler("hello")

	if ok {
		t.Error("Handler of wrong kind should not be valid")
	}

	if err == nil {
		t.Error("Handler of wrong kind should return error")
	}

	expectedError := fmt.Sprint("invalid handler type string, handler must be an func")
	if !strings.Contains(err.Error(), expectedError) {
		t.Errorf("Unexpected error message: \n\tActual: %v\n\tExpected: %v", err.Error(), expectedError)
	}
}

func TestNewingPanicsOnWrongInParameterCount(t *testing.T) {
	ok, err := validateHandler(HandlerWithWrongInParameterCount)

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
	ok, err := validateHandler(HandlerWithWrongOutParameterCount)

	if ok {
		t.Error("Handler with wrong parameter count should not be valid")
	}

	if err == nil {
		t.Fatal("Handler with wring parameter count should return error")
	}

	expectedError := fmt.Sprint("handler does not have 2 out parameters, instead it has 0")
	if !strings.Contains(err.Error(), expectedError) {
		t.Errorf("Unexpected error message: \n\tActual: %v\n\tExpected: %v", err.Error(), expectedError)
	}
}
