package model

import (
	"labix.org/v2/mgo/bson"
)

type ObjectId bson.ObjectId

func NewObjectId() ObjectId {
	id := bson.NewObjectId()
	s := id.String()

	return ObjectId(s)
}

func (id ObjectId) String() string {
	return string(id)
}
