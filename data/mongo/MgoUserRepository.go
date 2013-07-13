package mongo

import (
	"github.com/pjvds/httpcallback.io/model"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
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

func (r *MgoUserRepository) Get(id model.ObjectId) (*model.User, error) {
	query := r.database.C("Users").Find(bson.M{"_id": id})
	var result model.User
	err := query.One(&result)

	return &result, err
}
func (r *MgoUserRepository) GetByUsernameAndPasswordHash(username string, passwordHash string) (*model.User, error) {
	query := r.database.C("Users").Find(bson.M{"username": username, "passwordhash": passwordHash})
	var result *model.User
	if err := query.One(result); err != nil {
		return nil, err
	}

	if result == nil {
		return nil, nil
	}

	return result, nil
}

func (r *MgoUserRepository) GetByUsernameAndAuthToken(username string, authToken model.AuthenticationToken) (*model.UserAuthInfo, error) {
	query := r.database.C("Users").Find(bson.M{"username": username, "authtoken": authToken}).Select(bson.M{"_id": 1, "username": 1})
	var result bson.M
	if err := query.One(&result); err != nil {
		return nil, err
	}

	if result == nil {
		return nil, nil
	}

	return &model.UserAuthInfo{
		UserId:   model.ObjectId(result["_id"].(string)),
		Username: result["username"].(string),
	}, nil
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
