package model

import (
	"strings"
	"testing"
)

func TestParseObjectIdThrowsWhenLenghtIsNotCorrect(t *testing.T) {
	value := "foobar"
	_, err := ParseObjectId(value)

	if err == nil || !strings.Contains(err.Error(), "Invalid object id. String lenght") {
		t.Error("Expected ParseObjectId to fail with an error regarding invalid lenght")
	}
}

func TestParseObjectIdThrowsWhenNotValidHexString(t *testing.T) {
	value := NewObjectId().String()
	value = "G" + string(value[1:])

	_, err := ParseObjectId(value)

	if err == nil || !strings.Contains(err.Error(), "Invalid object id. Not a valid hexidecimal string") {
		t.Error("Expected ParseObjectId to fail with an error regarding invalid lenght")
	}
}
