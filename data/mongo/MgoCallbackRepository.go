package mongo

import (
	"github.com/pjvds/httpcallback.io/model"
	"labix.org/v2/mgo"
)

type MgoCallbackRepository struct {
	session  *MgoSession
	database *mgo.Database
}

func NewCallbackRepository(session *MgoSession) *MgoCallbackRepository {
	return &MgoCallbackRepository{
		session:  session,
		database: session.database,
	}
}

func (r *MgoCallbackRepository) Add(callback *model.Callback) error {
	return r.database.C("Callbacks").Insert(callback)
}

func (r *MgoCallbackRepository) List() ([]*model.Callback, error) {
	query := r.database.C("Callbacks").Find(nil)
	var result []*model.Callback
	err := query.All(&result)
	if result == nil {
		result = make([]*model.Callback, 0)
	}

	return result, err
}
