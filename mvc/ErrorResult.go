package mvc

import (
	"net/http"
)

type ErrorResultState struct {
	Error error
}

func ErrorResult(err error) *ErrorResultState {
	return &ErrorResultState{
		Error: err,
	}
}

func (r *ErrorResultState) WriteResponse(response http.ResponseWriter) {
	Log.Error("Writing Internal Server Error response because of: %v", r.Error.Error())

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusInternalServerError)
	response.Write([]byte("internal server error"))
}
