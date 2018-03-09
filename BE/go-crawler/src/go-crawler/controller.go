package main

import (
	"encoding/json"
	"log"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

/*
{
    "beginDate":
    {
        "year": 2018,
        "month": 10,
        "day": 9
    },
    "endDate":
    {
        "year": 2018,
        "month": 10,
        "day": 19
    }
}
*/

type datePost struct {
	Year  int        `json:"year"`
	Month time.Month `json:"month"`
	Day   int        `json:"day"`
}

// dateRange date range
type dateRangePost struct {
	BeginDate datePost `json:"beginDate"`
	EndDate   datePost `json:"endDate"`
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJSON(w, code, map[string]string{"error": msg})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err = w.Write(response)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}
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

// RangeReviewsEndPoint get range and respond with list of reviews in range
func RangeReviewsEndPoint(session *mgo.Session) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var dates dateRangePost
		if err := json.NewDecoder(r.Body).Decode(&dates); err != nil {
			respondWithError(w, http.StatusBadRequest, "bad request")
			log.Printf("bad request, error: %s,", err)

			return
		}

		beginDate := time.Date(
			dates.BeginDate.Year,
			dates.BeginDate.Month,
			dates.BeginDate.Day, 0, 0, 0, 0, time.UTC)
		endDate := time.Date(
			dates.EndDate.Year,
			dates.EndDate.Month,
			dates.EndDate.Day, 0, 0, 0, 0, time.UTC)

		log.Printf("range is from: %s to: %s", beginDate, endDate)

		// TODO: move to model methods
		c := session.DB(dbName).C(collectionName)

		var reviews []ReviewIndex
		err := c.Find(
			bson.M{
				"_time_stamp": bson.M{
					"$gt": beginDate,
					"$lt": endDate,
				},
			}).All(&reviews)

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
