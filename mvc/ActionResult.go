package mvc

import (
	"net/http"
)

// Represents a result of a controller method.
// The action can write the actual result to
// an http response stream.
//
// There are diffent types of ActionResults:
// - JsonResult
// - HttpStatusResult
// - etc...
type ActionResult interface {
	// Write the actual result to the http response stream.
	WriteResponse(http.ResponseWriter)
}
