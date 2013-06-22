package mongo

import (
	"github.com/pjvds/httpcallback.io/model"
	"labix.org/v2/mgo"
)

type MgoCallbackRepository struct {
	session  *mgo.Session
	database *mgo.Database
}

func NewCallbackRepository(url string, database string) (*MgoCallbackRepository, error) {
	session, err := mgo.Dial(url)
	if err != nil {
		return nil, err
	}
	err = session.Ping()
	if err != nil {
		return nil, err
	}
	db := session.DB(database)

	return &MgoCallbackRepository{
		session:  session,
		database: db,
	}, nil
}

func (r *MgoCallbackRepository) Add(callback *model.Callback) error {
	return r.database.C("Callbacks").Insert(callback)
}

func (r *MgoCallbackRepository) List() ([]*model.Callback, error) {
	query := r.database.C("Callbacks").Find(nil)
	var result []*model.Callback
	err := query.All(&result)

	return result, err
}
