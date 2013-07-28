package gocqrs

import (
	"github.com/pjvds/go-cqrs/sourcing"
	"github.com/pjvds/go-cqrs/storage"
	"github.com/pjvds/httpcallback.io/model"
)

type UserRepository struct {
	repository *storage.Repository
}

func (r *UserRepository) Add(user *model.User) error {
	source := sourcing.GetState(user)
	return r.repository.Add(source)
}

func (r *UserRepository) Get(sourceId model.ObjectId) (*model.User, error) {
	user := new(model.User)
	user.EventSource = sourcing.AttachNew(user)

	err := r.repository.Get(sourcing.EventSourceId(sourceId), user)
	return user, err
}
