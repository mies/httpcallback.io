package api

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"
)

func TestGetAuthorizationInfoFromRequestReadsFromHeader(t *testing.T) {
	username := "pjvds"
	token := "token"

	request, _ := http.NewRequest("POST", "http://httpcallback.io/users", nil)
	request.Header.Add("Authorization", fmt.Sprintf("HttpCallbackLogin username=\"%v\", token=\"%v\"", username, token))

	actualUsername, actualToken := GetAuthorizationInfoFromRequest(request)

	if username != actualUsername {
		t.Errorf("unexpected username: %v", actualUsername)
	}

	if token != actualToken {
		t.Errorf("unexpected token: %v", actualToken)
	}
}

func TestGetAuthorizationInfoFromRequestReadsFromQueryString(t *testing.T) {
	username := "pjvds"
	token := "token"

	s := fmt.Sprintf("http://f.com/users?auth_username=%v&auth_token=%v", url.QueryEscape(username), url.QueryEscape(token))
	request, _ := http.NewRequest("POST", s, nil)

	actualUsername, actualToken := GetAuthorizationInfoFromRequest(request)

	if username != actualUsername {
		t.Errorf("unexpected username: %v, expected: %v", actualUsername, username)
	}

	if token != actualToken {
		t.Errorf("unexpected token: %v, expected: %v", actualUsername, token)
	}
}
