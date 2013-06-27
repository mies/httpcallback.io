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

func TestTwoNewAuthTokensAreNotTheSame(t *testing.T) {
	newTokenA := NewAuthToken()
	newTokenB := NewAuthToken()
	if newTokenA != newTokenB {
		t.Error("New tokens should not be the same")
	}
}
