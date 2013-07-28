package gocqrs

import (
	"github.com/pjvds/go-cqrs/storage"
	"github.com/pjvds/go-cqrs/storage/eventstore"
	"github.com/pjvds/httpcallback.io/data"
)

type RepositoryFactory struct {
	repository *storage.Repository
}

func NewRepositoryFactory() *RepositoryFactory {
	backend, _ := eventstore.DailEventStore("http://localhost:2113", nil)

	return &RepositoryFactory{
		repository: storage.NewRepository(backend),
	}
}

func (r *RepositoryFactory) CreateUserRepository() data.UserRepository {
	return &UserRepository{
		repository: r.repository,
	}
}
