package data

import (
	"github.com/pjvds/httpcallback.io/model"
	"time"
)

type CallbackRepository interface {
	Add(*model.Callback) error
	List(userId model.ObjectId) ([]*model.Callback, error)

	AddAttemptToCallback(callbackId model.ObjectId, attempt *model.CallbackAttempt) error

	// Gets the next callback and bumps the NextAttemptTimeStamp
	// with the given duration.
	GetNextAndBumpNextAttemptTimeStamp(duration time.Duration) (*model.Callback, error)
}
