package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Review data
type Review struct {
	PostID int       `json:"post_id" bson:"post_id"`
	Date   time.Time `json:"date" bson:"date"`
	Author string    `json:"author" bson:"author"`
	Link   string    `json:"link" bson:"link"`
	Post   string    `json:"post" bson:"post"`
}

// ReviewIndex mongodb data struct
type ReviewIndex struct {
	ID        bson.ObjectId `json:"id" bson:"_id"`
	TimeStamp time.Time     `json:"time_stamp" bson:"_time_stamp"`
	Review
}

const (
	isDrop         = true
	dbName         = "test"
	collectionName = "review"
	dbURL          = "172.25.7.31:27017"
)

func storePost(session *mgo.Session) func(int, *goquery.Selection) {
	//return func(index int, item *goquery.Selection) {
	return func(index int, item *goquery.Selection) {
		//var scrapedReview review
		var scrapedReview ReviewIndex
		//title := item.Text()
		scrapedReview.PostID = index

		block := item.Find("blockquote").Text()
		scrapedReview.Post = strings.Join(strings.Split(strings.TrimSpace(block), "\n"), " ")

		scrapedReview.Author = item.Find("a").Text()

		dateString, _ := item.Find("time").Attr("datetime")
		scrapedReview.Date, _ = time.Parse(time.RFC3339, dateString)

		scrapedReview.Link, _ = item.Find("a").Attr("href")

		fmt.Printf("Post #%d: %s - %s - %s - %s\n", index, scrapedReview.Post, scrapedReview.Link, scrapedReview.Date, scrapedReview.Author)

		scrapedReview.ID = bson.NewObjectId()
		scrapedReview.TimeStamp = time.Now().UTC()

		c := session.DB(dbName).C(collectionName)

		errI := c.Insert(&scrapedReview)
		if errI != nil {
			panic(errI)
		}

		var reviews []ReviewIndex
		errF := c.Find(bson.M{}).All(&reviews)
		if errF != nil {
			panic(errF)
		}
		fmt.Printf("reviews: %v\n", reviews)

	}
}

func postScrape(session *mgo.Session) {
	doc, err := goquery.NewDocument("https://apps.shopify.com/omnisend#reviews-heading")
	if err != nil {
		log.Fatal(err)
	}

	doc.Find(".contents").Each(storePost(session))
}

func main() {
	// connect to mongo
	session, err := mgo.Dial(dbURL)
	if err != nil {
		panic(err)
	}

	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	defer session.Close()

	if isDrop {
		err = session.DB(dbName).DropDatabase()
		if err != nil {
			panic(err)
		}
	}

	postScrape(session)
}
