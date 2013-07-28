package model

import (
	"github.com/pjvds/go-cqrs/sourcing"
)

type ObjectId sourcing.EventSourceId

func NewObjectId() ObjectId {
	return ObjectId(sourcing.NewEventSourceId())
}

func (id *ObjectId) String() string {
	value := sourcing.EventSourceId(*id)
	return value.String()
}

func ParseObjectId(value string) (ObjectId, error) {
	id, err := sourcing.ParseEventSourceId(value)
	return ObjectId(id), err
}

func (id ObjectId) MarshalJSON() ([]byte, error) {
	value := sourcing.EventSourceId(id)
	return value.MarshalJSON()
}

func (id *ObjectId) UnmarshalJSON(b []byte) error {
	value := sourcing.EventSourceId(*id)
	err := value.UnmarshalJSON(b)
	if err != nil {
		return err
	}

	objectId := ObjectId(value)
	id = &objectId
	return nil
}
