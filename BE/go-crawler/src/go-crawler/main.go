package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

const (
	topList        = 5 // five most common
	shingleSize    = 3 // three-word phrases
	isDrop         = true
	dbName         = "test"
	collectionName = "review"
	bindPort       = "3000"
	dbURL          = "127.0.0.1:27017"
	scrapeURL      = "https://apps.shopify.com"
	firstURLSuffix = "/omnisend#reviews-heading" // start from the first page
	//firstURLSuffix = "/omnisend?page=130#reviews" // start from the 130 page
)

func main() {
	var mongo MgoSession
	mongo.Init(dbURL)

	log.Println("start scraping...")

	start := time.Now().UTC()
	scrapeSite(&mongo)
	timeScrape := time.Since(start)

	start = time.Now().UTC()
	printAllReviews(&mongo)
	timePrintAll := time.Since(start)

	cors := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedHeaders([]string{
			"X-Requested-With",
			"Origin",
			"Accept",
			"Authorization",
			"Content-type",
			"Access-Control-Allow-Origin",
		}),
		handlers.AllowedMethods([]string{
			http.MethodGet,
			http.MethodDelete,
			http.MethodPost,
			http.MethodPut,
			http.MethodHead,
			http.MethodOptions,
		}),
		handlers.AllowCredentials(),
		handlers.MaxAge(5),
	)

	r := mux.NewRouter()

	r.Use(cors)

	r.HandleFunc("/api/v1/reviews", AllReviewsEndPoint(&mongo)).Methods(http.MethodGet)
	r.HandleFunc("/api/v1/range", RangeReviewsEndPoint(&mongo)).Methods(http.MethodPost)
	r.HandleFunc("/api/v1/shingle", ShingleReviewsEndPoint(&mongo)).Methods(http.MethodPost)

	loggedRouter := handlers.LoggingHandler(os.Stdout, r)

	log.Printf("done scraping: %v", timeScrape)
	log.Printf("done printing: %v", timePrintAll)
	log.Printf("listening on port: %v", bindPort)

	if err := http.ListenAndServe(":"+bindPort, loggedRouter); err != nil {
		log.Fatal(err)
	}
}
