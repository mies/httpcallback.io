package memory

import (
	"github.com/pjvds/httpcallback.io/model"
	"sync"
)

type MemoryUserRepository struct {
	lock sync.RWMutex
	data []*model.User
}

func NewMemoryUserRepository() *MemoryUserRepository {
	return &MemoryUserRepository{
		data: make([]*model.User, 0),
	}
}

func (r *MemoryUserRepository) Add(user *model.User) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.data = append(r.data, user)

	Log.Debug("User added with id '%v' and username '%v'", user.Id, user.Username)
	Log.Debug("Total user count now: %v", len(r.data))
	return nil
}

func (r *MemoryUserRepository) List() ([]*model.User, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	return r.data, nil
}

func (r *MemoryUserRepository) Get(id model.ObjectId) (*model.User, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	length := len(r.data)
	Log.Debug("Getting user with id '%v' out of %v users", id, length)

	for i := 0; i < length; i++ {
		user := r.data[i]
		Log.Debug("Matching user %v with id '%v'", i+1, user.Id)

		if user.Id == id {
			return user, nil
		}
	}

	return nil, nil
}
