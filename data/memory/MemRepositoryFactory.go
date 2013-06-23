package memory

import (
	"github.com/pjvds/httpcallback.io/data"
)

type MemRepositoryFactory struct {
	users     *MemoryUserRepository
	callbacks *MemoryCallbackRepository
}

func NewMemRepositoryFactory() *MemRepositoryFactory {
	return &MemRepositoryFactory{
		users:     NewMemoryUserRepository(),
		callbacks: NewMemoryCallbackRepository(),
	}
}

func (f *MemRepositoryFactory) CreateUserRepository() data.UserRepository {
	return f.users
}

func (f *MemRepositoryFactory) CreateCallbackRepository() data.CallbackRepository {
	return f.callbacks
}
