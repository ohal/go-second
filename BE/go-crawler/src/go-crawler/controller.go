package main

import (
	//"log"
	"encoding/json"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJSON(w, code, map[string]string{"error": msg})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// AllReviewsEndPoint respond with list of reviews
func AllReviewsEndPoint(session *mgo.Session) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: move to model methods
		c := session.DB(dbName).C(collectionName)

		var reviews []ReviewIndex
		err := c.Find(bson.M{}).All(&reviews)
		if err != nil {
			panic(err)
		}

		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respondWithJSON(w, http.StatusOK, reviews)
	}
}
