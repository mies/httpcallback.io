package mongo

import (
	"github.com/pjvds/httpcallback.io/model"
	"labix.org/v2/mgo"
)

type MgoCallbackRepository struct {
	session *mgo.Session
}

func NewCallbackRepository(url string) (*MgoCallbackRepository, error) {
	session, err := mgo.Dial(url)
	if err != nil {
		return nil, err
	}

	return &MgoCallbackRepository{
		session: session,
	}, nil
}

func (r *MgoCallbackRepository) Add(callback *model.Callback) {
}

func (r *MgoCallbackRepository) List() []*model.Callback {
	return nil
}
