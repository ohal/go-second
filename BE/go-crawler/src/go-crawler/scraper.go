package main

import (
	//"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"gopkg.in/mgo.v2/bson"
)

// Review data
type Review struct {
	Date   time.Time `json:"date" bson:"date"`
	Author string    `json:"author" bson:"author"`
	Link   string    `json:"link" bson:"link"`
	Post   string    `json:"post" bson:"post"`
}

// ReviewIndex mongodb data
type ReviewIndex struct {
	ID        bson.ObjectId `json:"id" bson:"_id"`
	TimeStamp time.Time     `json:"time_stamp" bson:"_time_stamp"`
	Review
}

func scrapePage(mongo *MgoSession, url string) {
	log.Printf("url: %v\n", url)

	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find(".contents").Each(storePost(mongo))
}

func storePost(mongo *MgoSession) func(int, *goquery.Selection) {
	return func(index int, item *goquery.Selection) {
		var scrapedReview ReviewIndex

		scrapedReview.ID = bson.NewObjectId()

		block := item.Find("blockquote").Text()
		scrapedReview.Post = strings.Join(
			strings.Split(
				strings.TrimSpace(block), "\n"), " ")

		scrapedReview.Author = item.Find("a").Text()

		dateString, _ := item.Find("time").Attr("datetime")
		datetime, err := time.Parse(time.RFC3339, dateString)
		if err != nil {
			log.Fatal(err)
		}
		scrapedReview.TimeStamp = datetime
		scrapedReview.Date = datetime

		scrapedReview.Link, _ = item.Find("a").Attr("href")

		errI := mongo.Insert(&scrapedReview)
		if errI != nil {
			panic(errI)
		}
	}
}

func scrapeSite(mongo *MgoSession) {
	currentURL := scrapeURL + firstURLSuffix
	// scrape first page
	scrapePage(mongo, currentURL)
	// then check if we see next page
	var pageList sync.WaitGroup
	for { // scrape all pages until the end
		doc, err := goquery.NewDocument(currentURL)
		if err != nil {
			log.Fatal(err)
		}
		nextPageSuffix, _ := doc.Find(".next_page").Attr("href")
		// break loop if there is no page to scrape
		if nextPageSuffix == "" {
			break
		}
		// if there is a page to scrape then scrape it
		pageList.Add(1)
		currentURL = scrapeURL + nextPageSuffix
		go func() {
			defer pageList.Done()
			scrapePage(mongo, currentURL)
		}()
	}
	pageList.Wait()
}
