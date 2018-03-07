package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	isDrop         = true
	dbName         = "test"
	collectionName = "review"
	dbURL          = "172.25.7.31:27017"
	scrapeURL      = "https://apps.shopify.com"
	//firstURLSuffix = "/omnisend#reviews-heading" // start from the first page
	firstURLSuffix = "/omnisend?page=130#reviews" // start from the 130 page
)

func main() {
	var mongo MgoSession
	session := mongo.Init(dbURL)

	scrapeSite(session)
	printAllReviews(session)

	r := mux.NewRouter()
	r.HandleFunc("/api/v1/reviews", AllReviewsEndPoint(session)).Methods("GET")
	if err := http.ListenAndServe(":3000", r); err != nil {
		log.Fatal(err)
	}
}
