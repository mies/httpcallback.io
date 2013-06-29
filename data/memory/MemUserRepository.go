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
	return r.get(func(user *model.User) bool {
		return user.Id == id
	})
}

func (r *MemoryUserRepository) GetByAuth(username string, authToken model.AuthenticationToken) (*model.UserAuthInfo, error) {
	user, err := r.get(func(user *model.User) bool {
		return user.Username == username && user.AuthToken == authToken
	})

	if user == nil || err != nil {
		return nil, err
	}

	return &model.UserAuthInfo{
		UserId:   user.Id,
		Username: user.Username,
	}, nil
}

func (r *MemoryUserRepository) get(predicate func(*model.User) bool) (*model.User, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	length := len(r.data)
	for i := 0; i < length; i++ {
		user := r.data[i]
		isMatch := predicate(user)

		if isMatch {
			return user, nil
		}
	}

	return nil, nil
}
