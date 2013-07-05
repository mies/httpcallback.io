package mvc

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// A simple Json Document type which you can
// use in the following way:
//
// doc := JsonDocument {
//   "Foo": "Bar",
//   "Baz": 42,
// }
type JsonDocument map[string]interface{}

type JsonResponse struct {
	data []byte
}

// Marshals the result object to a json string
// that will be written to the body of the response.
// This function panics if the result could not be
// marshalled.
func JsonResult(result interface{}) *JsonResponse {
	data, err := json.Marshal(result)
	if err != nil {
		message := fmt.Sprintf("Unable to marshal object (%+v) to json: %v", result, err.Error())
		panic(message)
	} else {
		return &JsonResponse{
			data: data,
		}
	}
}

// Writes the json to the HttpResponse stream. It also sets
// the Content-Type header value to application/json.
func (j *JsonResponse) WriteResponse(response http.ResponseWriter) {
	response.Header().Set("Content-Type", "application/json")
	response.Write(j.data)
}
