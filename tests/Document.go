package tests

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Document map[string]interface{}

func GetBodyAsDocument(response *http.Response) (Document, error) {
	var document Document
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return document, err
	}

	err = json.Unmarshal(data, &document)
	if err != nil {
	}

	return document, err
}

func (doc Document) ToJson() []byte {
	data, err := json.Marshal(doc)
	if err != nil {
		panic(err)
	}
	return data
}
