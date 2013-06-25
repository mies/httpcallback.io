package security

import (
	"testing"
)

func TestNewAuthToken(t *testing.T) {
	var emptyToken AuthenticationToken
	newToken := NewAuthToken()
	if newToken == emptyToken {
		t.Error("New token should not be empty")
	}
}
