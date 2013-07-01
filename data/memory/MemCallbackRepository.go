package memory

import (
	"errors"
	"github.com/pjvds/httpcallback.io/model"
	"sync"
	"time"
)

type MemoryCallbackRepository struct {
	lock sync.RWMutex
	data []*model.Callback
}

func NewMemoryCallbackRepository() *MemoryCallbackRepository {
	return &MemoryCallbackRepository{
		data: make([]*model.Callback, 0),
	}
}

func (r *MemoryCallbackRepository) Add(callback *model.Callback) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.data = append(r.data, callback)
	return nil
}

func (r *MemoryCallbackRepository) AddAttemptToCallback(callbackId model.ObjectId, attempt *model.CallbackAttempt) error {
	r.lock.RLock()
	defer r.lock.RUnlock()

	for _, c := range r.data {
		if c.Id == callbackId {
			r.lock.Lock()
			defer r.lock.Unlock()

			c.Attempts = append(c.Attempts, attempt)
			c.AttemptCount = len(c.Attempts)
			return nil
		}
	}

	return errors.New("not found")
}

func (r *MemoryCallbackRepository) GetNextAndBumpNextAttemptTimeStamp(duration time.Duration) (*model.Callback, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	now := time.Now()
	for _, c := range r.data {
		if !c.Finished && c.NextAttemptTimeStamp.After(now) {
			r.lock.Lock()
			defer r.lock.Unlock()

			c.NextAttemptTimeStamp = now.Add(duration)
			return c, nil
		}
	}

	return nil, nil
}

func (r *MemoryCallbackRepository) List(userId model.ObjectId) ([]*model.Callback, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()
	result := make([]*model.Callback, 0)

	for _, c := range r.data {
		if c.UserId == userId {
			result = append(result, c)
		}
	}

	return result, nil
}
