package security

import (
	"github.com/pjvds/httpcallback.io/model"
	"testing"
)

func TestNewAuthToken(t *testing.T) {
	var emptyToken model.AuthenticationToken
	newToken := NewAuthToken()
	if newToken == emptyToken {
		t.Error("New token should not be empty")
	}
}

func TestTwoNewAuthTokensAreNotTheSame(t *testing.T) {
	newTokenA := NewAuthToken()
	newTokenB := NewAuthToken()
	if newTokenA == newTokenB {
		t.Errorf("New tokens should not be the same:\n\ttoken a: %v\n\ttoken b: %v",
			newTokenA, newTokenB)
	}
}
