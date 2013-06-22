package mongo

import (
	"labix.org/v2/mgo"
)

type MgoCallbackRepository struct {
	db mgo.Database
}

func NewCallbackRepository(url string) (*MgoCallbackRepository, error) {
	session, err := mgo.Dail(url)
	if err != nil {
		return nil, err
	}

	return &MgoCallbackRepository{
		db: session,
	}, nil
}

func (r *MgoCallbackRepository) Add(callback *model.Callback) {
}

func (r *MgoCallbackRepository) List() []*model.Callback {
	return nil
}
