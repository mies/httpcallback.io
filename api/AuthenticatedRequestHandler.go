package api

import (
	"net/http"
)

// Represent an handher that only handles AuthenticatedRequests. Requests that
// are not successfully authenticated will never make it to this handler.
type AuthenticatedRequestHandler func(http.ResponseWriter, *AuthenticatedRequest)
