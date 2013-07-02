package mongo

import (
	"labix.org/v2/mgo"
)

type MgoSession struct {
	session  *mgo.Session
	database *mgo.Database
}

func Open(url string, database string) (*MgoSession, error) {
	Log.Info("Dailing Mongo at %v", url)

	session, err := mgo.Dial(url)
	if err != nil {
		Log.Error("Dailing Mongo failed: %v", err.Error())
		return nil, err
	}

	Log.Debug("Dailed Mongo successfully, start ping")
	err = session.Ping()
	if err != nil {
		Log.Error("Pinging Mongo failed: %v", err.Error())
		return nil, err
	}

	Log.Debug("Switching to database %v", database)
	db := session.DB(database)

	return &MgoSession{
		session:  session,
		database: db,
	}, nil
}
