package mongo

import (
	"github.com/pjvds/httpcallback.io/data"
)

type MgoRepositoryFactory struct {
	users     *MgoUserRepository
	callbacks *MgoCallbackRepository
}

func NewMgoRepositoryFactory(session *MgoSession) *MgoRepositoryFactory {
	return &MgoRepositoryFactory{
		users:     NewUserRepository(session),
		callbacks: NewCallbackRepository(session),
	}
}

func (f *MgoRepositoryFactory) CreateUserRepository() data.UserRepository {
	return f.users
}

func (f *MgoRepositoryFactory) CreateCallbackRepository() data.CallbackRepository {
	return f.callbacks
}
