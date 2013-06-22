package api

import (
	"encoding/json"
)

type Response interface {
	ToJson() []byte
}

type ResponseData map[string]interface{}

func (r ResponseData) ToJson() []byte {
	bytes, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}
	return bytes
}
