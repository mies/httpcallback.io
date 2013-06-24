package security

import (
	"strings"
	"testing"
	"time"
)

func TestHashPasswordNotContainsInputs(t *testing.T) {
	var username = "pjvds"
	var password = "foobar"
	var createdAt, _ = time.Parse(time.RFC822, "Wed, 02 Oct 2002 13:00:00 GMT")

	hash := HashPassword(username, password, createdAt)

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
	var createdAt, _ = time.Parse(time.RFC822, "Wed, 02 Oct 2002 13:00:00 GMT")

	hashA := HashPassword(usernameA, password, createdAt)
	hashB := HashPassword(usernameB, password, createdAt)

	if hashA == hashB {
		t.Error("Hashes should not be the same if username is different")
	}
}

func TestHashPasswordChangesForPassword(t *testing.T) {
	var username = "pjvds"
	var passwordA = "foobar"
	var passwordB = "baz"
	var createdAt, _ = time.Parse(time.RFC822, "Wed, 02 Oct 2002 13:00:00 GMT")

	hashA := HashPassword(username, passwordA, createdAt)
	hashB := HashPassword(username, passwordB, createdAt)

	if hashA == hashB {
		t.Error("Hashes should not be the same if password is different")
	}
}

func TestHashPasswordChangesForTimestamp(t *testing.T) {
	var username = "pjvds"
	var password = "foobar"
	var createdAtA = time.Now()
	var createdAtB = createdAtA.Add(5 * time.Hour)

	hashA := HashPassword(username, password, createdAtA)
	hashB := HashPassword(username, password, createdAtB)

	if hashA == hashB {
		t.Error("Hashes should not be the same if timestamp is different")
	}
}
