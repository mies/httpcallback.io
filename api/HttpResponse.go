package api

import (
	"net/http"
)

type ActionResult interface {
	WriteResponse(http.ResponseWriter)
}
