package memory

import (
	"github.com/pjvds/httpcallback.io/model"
)

type MemoryUserRepository struct {
	data []*model.User
}

func NewMemoryUserRepository() *MemoryUserRepository {
	return &MemoryUserRepository{
		data: make([]*model.User, 0),
	}
}

func (r *MemoryUserRepository) Add(user *model.User) error {
	r.data = append(r.data, user)
	return nil
}

func (r *MemoryUserRepository) List() ([]*model.User, error) {
	return r.data, nil
}
