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
	url     string
}

// Init init DB connection
func (s *MgoSession) Init(url string) {
	var err error

	s.url = url
	// connect to mongo
	s.session, err = mgo.Dial(url)
	if err != nil {
		panic(err)
	}

	s.session.SetMode(mgo.Monotonic, true)
	defer s.session.Close()

	if isDrop {
		err = s.session.DB(dbName).DropDatabase()
		if err != nil {
			panic(err)
		}
	}
}

// GetSession get session
func (s *MgoSession) GetSession() *mgo.Session {
	return s.session.Copy()
}

// CloseSession close session
func (s *MgoSession) CloseSession() {
	s.session.Close()
}

// FindAll find all reviews
func (s *MgoSession) FindAll() ([]ReviewIndex, error) {
	var err error

	// connect to mongo
	s.session, err = mgo.Dial(s.url)
	if err != nil {
		panic(err)
	}

	defer s.session.Close()

	s.session.SetMode(mgo.Monotonic, true)

	c := s.session.DB(dbName).C(collectionName)

	var reviews []ReviewIndex
	err = c.Find(bson.M{}).All(&reviews)
	if err != nil {
		panic(err)
	}
	return reviews, err
}

// FindRange find all reviews from/to ranged
func (s *MgoSession) FindRange(beginDate, endDate time.Time) ([]ReviewIndex, error) {
	var err error

	// connect to mongo
	s.session, err = mgo.Dial(s.url)
	if err != nil {
		panic(err)
	}

	defer s.session.Close()

	s.session.SetMode(mgo.Monotonic, true)

	c := s.session.DB(dbName).C(collectionName)

	var reviews []ReviewIndex
	err = c.Find(
		bson.M{
			"_time_stamp": bson.M{
				"$gt": beginDate,
				"$lt": endDate,
			},
		}).All(&reviews)
	if err != nil {
		panic(err)
	}
	return reviews, err
}

// Insert insert to mgo
func (s *MgoSession) Insert(scrapedReview *ReviewIndex) error {
	var err error

	// connect to mongo
	s.session, err = mgo.Dial(s.url)
	if err != nil {
		panic(err)
	}

	defer s.session.Close()

	s.session.SetMode(mgo.Monotonic, true)

	c := s.session.DB(dbName).C(collectionName)

	err = c.Insert(scrapedReview)
	if err != nil {
		panic(err)
	}
	return err
}

func printAllReviews(mongo *MgoSession) {
	reviews, err := mongo.FindAll()
	if err != nil {
		panic(err)
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
	var postList sync.WaitGroup
	postList.Add(len(reviews))
	for _, r := range reviews {
		go func(r ReviewIndex) {
			posts <- post{r.ID, r.TimeStamp, r.Date, r.Author, r.Link, r.Post}
		}(r)
	}

	go func() {
		for post := range posts {
			log.Printf("post: %v\n", post)
			postList.Done()
		}
	}()
	postList.Wait()
}
