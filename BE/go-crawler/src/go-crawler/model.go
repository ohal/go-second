package main

import (
	"log"
	"sync"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// MgoSession public
type MgoSession struct {
	session *mgo.Session
}

// Init init DB connection
func (s *MgoSession) Init(url string) *mgo.Session {
	var err error
	// connect to mongo
	s.session, err = mgo.Dial(url)
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

// FindAll find all reviews
func (s *MgoSession) FindAll() ([]ReviewIndex, error) {
	c := s.session.DB(dbName).C(collectionName)

	var reviews []ReviewIndex
	err := c.Find(bson.M{}).All(&reviews)
	if err != nil {
		panic(err)
	}
	return reviews, err
}

func printAllReviews(session *mgo.Session) {
	c := session.DB(dbName).C(collectionName)

	var reviews []ReviewIndex
	errF := c.Find(bson.M{}).All(&reviews)
	if errF != nil {
		panic(errF)
	}

	type post struct {
		ID        bson.ObjectId `json:"id"`
		TimeStamp time.Time     `json:"time_stamp"`
		Date      time.Time     `json:"date"`
		Author    string        `json:"author"`
		Link      string        `json:"link"`
		Post      string        `json:"post"`
	}
	posts := make(chan post)
	var phrasesList sync.WaitGroup
	phrasesList.Add(len(reviews))
	for _, r := range reviews {
		go func(r ReviewIndex) {
			defer phrasesList.Done()
			posts <- post{r.ID, r.TimeStamp, r.Date, r.Author, r.Link, r.Post}
		}(r)
	}

	go func() {
		for post := range posts {
			log.Printf("post: %v\n", post)
		}
	}()
	phrasesList.Wait()
}
