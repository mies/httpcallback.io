package api

import (
	"net/http"
)

// Represent an handher that only handles AuthenticatedRequests. Request that
// are not successfully authenticated will never make it to this handler.
type AuthenticatedRequestHandler func(http.ResponseWriter, *AuthenticatedRequest)
