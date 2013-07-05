package model

import (
	"encoding/hex"
	"errors"
	"fmt"
	"labix.org/v2/mgo/bson"
)

type ObjectId string

func NewObjectId() ObjectId {
	id := bson.NewObjectId()
	s := id.Hex()
	return ObjectId(s)
}

func (id ObjectId) String() string {
	return string(id)
}

func ParseObjectId(value string) (ObjectId, error) {
	var id ObjectId
	if len(value) != 24 {
		return id, errors.New(fmt.Sprintf("Invalid object id. String lenght is %v while it must be %v", len(value), 24))
	}

	_, err := hex.DecodeString(value)
	if err != nil {
		return id, errors.New(fmt.Sprintf("Invalid object id. Not a valid hexidecimal string: %v", err.Error()))
	}

	id = ObjectId(bson.ObjectIdHex(value).Hex())
	return id, nil
}
