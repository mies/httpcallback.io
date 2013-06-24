package model

import (
	"labix.org/v2/mgo/bson"
)

type ObjectId string

func NewObjectId() ObjectId {
	id := bson.NewObjectId()
	s := id.String()

	return ObjectId(s)
}
