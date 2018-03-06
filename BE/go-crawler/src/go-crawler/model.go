package main

import (
	"gopkg.in/mgo.v2"
)

// MgoSession public
type MgoSession struct {
	session *mgo.Session
}

// Init init DB connection
func (s *MgoSession) Init() *mgo.Session {
	var err error
	// connect to mongo
	s.session, err = mgo.Dial(dbURL)
	if err != nil {
		panic(err)
	}

	defer s.session.Close()

	s.session.SetMode(mgo.Monotonic, true)
	defer s.session.Close()

	if isDrop {
		err = s.session.DB(dbName).DropDatabase()
		if err != nil {
			panic(err)
		}
	}

	return s.session.Copy()
}

// GetSession get session
func (s *MgoSession) GetSession() *mgo.Session {
	return s.session.Copy()
}

// CloseSession close session
func (s *MgoSession) CloseSession() {
	s.session.Close()
	return
}
