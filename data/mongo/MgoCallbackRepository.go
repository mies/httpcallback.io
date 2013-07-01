package mongo

import (
	"github.com/pjvds/httpcallback.io/model"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"time"
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

func (r *MgoCallbackRepository) GetNextAndBumpNextAttemptTimeStamp(duration time.Duration) (*model.Callback, error) {
	now := time.Now()

	change := mgo.Change{
		Update:    bson.M{"$set": bson.M{"nextattempttimestamp": now.Add(duration)}},
		ReturnNew: true,
	}

	var result model.Callback
	_, err := r.database.C("Callbacks").Find(bson.M{"nextattempttimestamp": bson.M{"$lt": now}}).Apply(change, &result)
	if err != nil {
		if err.Error() == mgo.ErrNotFound.Error() {
			// mgo returns an ErrNotFound for atomic findAndModify operation
			// we just return nil in that case.
			return nil, nil
		}
		return nil, err
	}
	return &result, nil
}

func (r *MgoCallbackRepository) AddAttemptToCallback(callbackId model.ObjectId, attempt *model.CallbackAttempt) error {
	return r.database.C("Callbacks").UpdateId(callbackId, bson.M{
		"$push": bson.M{"attempts": attempt},
		"$inc":  bson.M{"attemptcount": 1},
		"$set":  bson.M{"finished": attempt.Success},
	})
}

func (r *MgoCallbackRepository) List(userId model.ObjectId) ([]*model.Callback, error) {
	query := r.database.C("Callbacks").Find(bson.M{"userid": userId})
	var result []*model.Callback
	err := query.All(&result)
	if result == nil {
		result = make([]*model.Callback, 0)
	}

	return result, err
}
