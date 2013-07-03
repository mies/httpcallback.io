package mvc

import (
	"encoding/json"
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

func JsonResult(result interface{}) (*JsonResponse, error) {
	data, err := json.Marshal(result)
	if err != nil {
		Log.Error("Unable to marshal object (%+v) to json: %s", result, err.Error())
		return nil, err
	} else {
		return &JsonResponse{
			data: data,
		}, nil
	}
}

// Writes the json to the HttpResponse stream. It also sets
// the Content-Type header value to application/json.
func (j *JsonResponse) WriteResponse(response http.ResponseWriter) {
	response.Header().Set("Content-Type", "application/json")
	response.Write(j.data)
}
