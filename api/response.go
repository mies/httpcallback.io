package api

import (
	"encoding/json"
)

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func (r Response) ToJson() []byte {
	bytes, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}
	return bytes
}
