package api

import (
	"github.com/pjvds/httpcallback.io/data"
	"github.com/pjvds/httpcallback.io/model"
	"net/http"
	//"regexp"
)

type AuthenticatedRequestHandler func(http.ResponseWriter, *AuthenticatedRequest)

type AuthenticatedRequest struct {
	*http.Request

	// The user that made the request.
	UserId   model.ObjectId
	Username string
}

type AuthenticationHandler struct {
	UserRepository         data.UserRepository
	Handler                AuthenticatedRequestHandler
	NotFoundOnUnauthorized bool
}

func (h *AuthenticationHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	username := request.URL.Query().Get("auth_username")
	token := request.URL.Query().Get("auth_token")

	if username == "" || token == "" {
		Log.Warning("No auth info supplied in request")

		if h.NotFoundOnUnauthorized {
			response.WriteHeader(http.StatusNotFound)
		} else {
			response.WriteHeader(http.StatusUnauthorized)
		}
		return
	}

	user, err := h.UserRepository.GetByAuth(username, model.AuthenticationToken(token))

	if err != nil {
		Log.Error("Unable to authenticate request: %v", err.Error())
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	if user == nil {
		Log.Warning("Not authenticated! User not found by username=%v and token=<hidden>", username)

		if h.NotFoundOnUnauthorized {
			response.WriteHeader(http.StatusNotFound)
		} else {
			response.WriteHeader(http.StatusUnauthorized)
		}
		return
	}

	authRequest := &AuthenticatedRequest{
		Request:  request,
		UserId:   user.UserId,
		Username: user.Username,
	}

	h.Handler(response, authRequest)
}
