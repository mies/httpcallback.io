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

func (r *MemoryCallbackRepository) List(userId model.ObjectId) ([]*model.Callback, error) {
	result := make([]*model.Callback, 0)

	for _, c := range r.data {
		if c.UserId == userId {
			result = append(result, c)
		}
	}

	return result, nil
}
