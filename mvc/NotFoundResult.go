package mvc

import (
	"net/http"
)

type NotFoundResultState struct {
	message string
}

func NotFoundResult(message string) *NotFoundResultState {
	return &NotFoundResultState{
		message: message,
	}
}

func (r *NotFoundResultState) WriteResponse(response http.ResponseWriter) {
	jsonResult := JsonResult(&JsonDocument{
		"message": r.message,
	})

	response.WriteHeader(http.StatusNotFound)
	jsonResult.WriteResponse(response)
}
