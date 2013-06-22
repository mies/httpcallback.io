package mongo

import (
	"github.com/pjvds/httpcallback.io/model"
	"labix.org/v2/mgo"
)

type MgoUserRepository struct {
	session  *MgoSession
	database *mgo.Database
}

func NewUserRepository(session *MgoSession) *MgoUserRepository {
	return &MgoUserRepository{
		session:  session,
		database: session.database,
	}
}

func (r *MgoUserRepository) Add(user *model.User) error {
	return r.database.C("Users").Insert(user)
}

func (r *MgoUserRepository) List() ([]*model.User, error) {
	query := r.database.C("Users").Find(nil)
	var result []*model.User
	err := query.All(&result)
	if result == nil {
		result = make([]*model.User, 0)
	}

	return result, err
}
