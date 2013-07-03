package mvc

import (
	"net/http"
)

// Represents an handler that authenticates incoming requests. Whenver an
// incoming request is successfully authenticated the
// AuthenticationHandler.Handler is called.
type AuthenticationHandler struct {
	authenticator Authenticator

	// The handler that will handle all requests that are successfully authenticated.
	Handler AuthenticatedRequestHandler

	// Indicates whether unauthorized access should result in StatusUnauthorized
	// or StatusNotFound result. This can be hulpful if private urls should not
	// be discoverable for unauthorized users. Default is false.
	NotFoundOnUnauthorized bool
}

// Serve an HTTP request. It will validate the authentication information and
// call the AuthenticationHandler.Handler whenever the request was successfully
// authenticated. If authentication fails a StatusUnauthorized or StatusNotFound
// is returned to the client. This depends on the
// AuthenticationHandler.NotFoundOnUnauthorized flag.
//
// Whenever the Authenticator returns an error it will be logged and a
// StatusInternalServerError will be send to the client.
func (h *AuthenticationHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	username, token := GetAuthorizationInfoFromRequest(request)

	// Validate whether we received authentication information.
	if username == "" || token == "" {
		Log.Warning("No auth info supplied in request")

		if h.NotFoundOnUnauthorized {
			response.WriteHeader(http.StatusNotFound)
			response.Write([]byte("Not found"))
		} else {
			response.WriteHeader(http.StatusUnauthorized)
			response.Write([]byte("Not authorized"))
		}
		return
	}

	// Do actual authentication.
	user, err := h.authenticator.Authenticate(username, token)

	// Did the authentication process fail?
	if err != nil {
		Log.Error("Unable to authenticate request: %v", err.Error())
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Did the authenticator gave us a user?
	if user == nil {
		Log.Warning("Not authenticated! User not found by username=%v and token=<hidden>", username)

		if h.NotFoundOnUnauthorized {
			response.WriteHeader(http.StatusNotFound)
		} else {
			response.WriteHeader(http.StatusUnauthorized)
		}
		return
	}

	// Request was authenticated successfully, set information and call
	// actual handler.
	authRequest := NewAuthenticatedRequest(request, user.UserId, user.Username)
	h.Handler(response, authRequest)
}
