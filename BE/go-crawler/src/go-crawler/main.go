package main

import (
	"log"
)

const (
	isDrop         = true
	dbName         = "test"
	collectionName = "review"
	dbURL          = "172.25.7.31:27017"
)

func main() {
	var mongo MgoSession
	session := mongo.Init()
	log.Printf("session: %v\n", session)

	postScrape(session)
}
