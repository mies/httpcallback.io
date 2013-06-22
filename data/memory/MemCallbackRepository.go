package memory

import (
	"github.com/pjvds/httpcallback.io/model"
)

type MemoryCallbackRepository struct {
	data []*model.Callback
}

func NewMemoryCallbackRepository() *MemoryCallbackRepository {
	return &MemoryCallbackRepository{
		data: make([]*model.Callback, 0),
	}
}

func (r *MemoryCallbackRepository) Add(callback *model.Callback) error {
	r.data = append(r.data, callback)
	return nil
}

func (r *MemoryCallbackRepository) List() ([]*model.Callback, error) {
	return r.data, nil
}
