package mongo

import (
	"labix.org/v2/mgo"
)

type MgoSession struct {
	session  *mgo.Session
	database *mgo.Database
}

func Open(url string, database string) (*MgoSession, error) {
	session, err := mgo.Dial(url)
	if err != nil {
		return nil, err
	}
	err = session.Ping()
	if err != nil {
		return nil, err
	}
	db := session.DB(database)

	return &MgoSession{
		session:  session,
		database: db,
	}, nil
}
