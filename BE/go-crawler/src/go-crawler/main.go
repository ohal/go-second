package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

const (
	isDrop         = true
	dbName         = "test"
	collectionName = "review"
	dbURL          = "127.0.0.1:27017"
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
	loggedRouter := handlers.LoggingHandler(os.Stdout, r)
	if err := http.ListenAndServe(":3000", loggedRouter); err != nil {
		log.Fatal(err)
	}
}
