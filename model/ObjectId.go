package model

import (
	"labix.org/v2/mgo/bson"
)

type ObjectId bson.ObjectId

func NewObjectId() ObjectId {
	return ObjectId(bson.NewObjectId())
}
