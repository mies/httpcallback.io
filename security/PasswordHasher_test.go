package security

import (
	"strings"
	"testing"
)

func TestHashPasswordNotContainsInputs(t *testing.T) {
	var username = "pjvds"
	var password = "foobar"

	hash := HashPassword(username, password)

	if strings.Contains(hash, username) {
		t.Error("Hash should not contains username")
	}

	if strings.Contains(hash, password) {
		t.Error("Hash should not contains username")
	}
}

func TestHashPasswordChangesForUsername(t *testing.T) {
	var usernameA = "pjvds"
	var usernameB = "BAZ"
	var password = "foobar"

	hashA := HashPassword(usernameA, password)
	hashB := HashPassword(usernameB, password)

	if hashA == hashB {
		t.Error("Hashes should not be the same if username is different")
	}
}

func TestHashPasswordChangesForPassword(t *testing.T) {
	var username = "pjvds"
	var passwordA = "foobar"
	var passwordB = "baz"

	hashA := HashPassword(username, passwordA)
	hashB := HashPassword(username, passwordB)

	if hashA == hashB {
		t.Error("Hashes should not be the same if password is different")
	}
}
