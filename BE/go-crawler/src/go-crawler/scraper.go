package main

import (
	//"fmt"
	"log"
	"strings"
	"sync"
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

func storePost(session *mgo.Session) func(int, *goquery.Selection) {
	return func(index int, item *goquery.Selection) {
		var scrapedReview ReviewIndex

		scrapedReview.PostID = index

		block := item.Find("blockquote").Text()
		scrapedReview.Post = strings.Join(strings.Split(strings.TrimSpace(block), "\n"), " ")

		scrapedReview.Author = item.Find("a").Text()

		dateString, _ := item.Find("time").Attr("datetime")
		scrapedReview.Date, _ = time.Parse(time.RFC3339, dateString)

		scrapedReview.Link, _ = item.Find("a").Attr("href")

		//log.Printf("Post #%d: %s - %s - %s - %s\n", index, scrapedReview.Post, scrapedReview.Link, scrapedReview.Date, scrapedReview.Author)

		scrapedReview.ID = bson.NewObjectId()
		scrapedReview.TimeStamp = time.Now().UTC()

		c := session.DB(dbName).C(collectionName)

		errI := c.Insert(&scrapedReview)
		if errI != nil {
			panic(errI)
		}
	}
}

func printAllReviews(session *mgo.Session) {
	c := session.DB(dbName).C(collectionName)

	var reviews []ReviewIndex
	errF := c.Find(bson.M{}).All(&reviews)
	if errF != nil {
		panic(errF)
	}
	log.Printf("reviews: %v\n", reviews)

	phrases := make(chan string, len(reviews))
	var postList sync.WaitGroup
	postList.Add(len(reviews))
	go func() {
		for _, r := range reviews {
			phrases <- r.Post
		}
		close(phrases)
	}()

	for p := range phrases {
		go func() {
			log.Printf("p: %v\n", p)
			postList.Done()
		}()
	}

	postList.Wait()
}

func postScrape(session *mgo.Session) {
	doc, err := goquery.NewDocument("https://apps.shopify.com/omnisend#reviews-heading")
	if err != nil {
		log.Fatal(err)
	}

	doc.Find(".contents").Each(storePost(session))

	printAllReviews(session)
}
