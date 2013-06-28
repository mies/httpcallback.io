package api

// import (
// 	"github.com/pjvds/httpcallback.io/model"
// 	"net/http"
// )

// type AuthenticatedRequest struct {
// 	http.Request

// 	// The user that made the request.
// 	User *model.User
// }

// type AuthenticationHandler struct {
// 	Predicate AuthenicateRequestPredicate
// 	Handler   http.Handler
// }

// type AuthenicateRequestPredicate func(http.Request) (model.User, bool, error)

// func (h *AuthenticationHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
// 	user, ok, err := h.Predicate(request)
// 	if err != nil {
// 		Log.Error("Unable to authenticate request: %v", err.Error())
// 		response.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}

// 	if !ok {
// 		Log.Info("Not authenticated!")
// 		response.WriteHeader(http.StatusUnauthorized)
// 		return
// 	}

// 	h.Handler.ServeHTTP(response, request)
// }
