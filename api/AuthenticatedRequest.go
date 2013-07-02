package api

import (
	"github.com/pjvds/httpcallback.io/model"
	"net/http"
)

// A AuthenticatedRequest represents an HTTP request received by the server
// that has been authenticated against. Only request that are succesfully
// authenticated will make it into an AuthenticatedRequest.
type AuthenticatedRequest struct {
	// The actual HTTP request. This is always set.
	*http.Request

	// The user that made the request. This property is always set, thus never
	// nil or empty.
	UserId model.ObjectId
	// The username of the user that has been authenticated. This property
	// is always set, thus never nil or empty.
	Username string
}

// Creates a new AuthenticatedRequest. Panics when request is nil, or username is empty.
func NewAuthenticatedRequest(request *http.Request, userId model.ObjectId, username string) *AuthenticatedRequest {
	if request == nil {
		panic("request cannot be nil")
	}
	if username == "" {
		panic("username cannot be empty")
	}

	return &AuthenticatedRequest{
		Request:  request,
		UserId:   userId,
		Username: username,
	}
}
