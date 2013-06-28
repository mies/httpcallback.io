package api

import (
	"fmt"
	"strings"
	"testing"
)

func HandlerWithWrongParameterCount() {
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

func TestNewingPanicsOnWrongArgumentCount(t *testing.T) {
	ok, err := validateHandler(HandlerWithWrongParameterCount)

	if ok {
		t.Error("Handler with wrong parameter count should not be valid")
	}

	if err == nil {
		t.Error("Handler with wring parameter count should return error")
	}

	expectedError := fmt.Sprint("handler does not have 2 arguments, instead it has 0")
	if !strings.Contains(err.Error(), expectedError) {
		t.Errorf("Unexpected error message: \n\tActual: %v\n\tExpected: %v", err.Error(), expectedError)
	}
}
