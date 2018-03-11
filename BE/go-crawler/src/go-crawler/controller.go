package main

import (
	"encoding/json"
	"log"
	"sync"
	"time"

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
func AllReviewsEndPoint(mongo *MgoSession) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reviews, err := mongo.FindAll()
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respondWithJSON(w, http.StatusOK, reviews)
	}
}

// RangeReviewsEndPoint get range and respond with list of reviews in range
func RangeReviewsEndPoint(mongo *MgoSession) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var dates dateRangePost
		if err := json.NewDecoder(r.Body).Decode(&dates); err != nil {
			respondWithError(w, http.StatusBadRequest, "bad request")
			log.Printf("bad request, error: %s,", err)

			return
		}

		beginDate, endDate := getRange(dates)
		log.Printf("range is from: %s to: %s", beginDate, endDate)

		reviews, err := mongo.FindRange(beginDate, endDate)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respondWithJSON(w, http.StatusOK, reviews)
	}
}

// ShingleReviewsEndPoint get range and respond with list of reviews in range
func ShingleReviewsEndPoint(mongo *MgoSession) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var dates dateRangePost
		if err := json.NewDecoder(r.Body).Decode(&dates); err != nil {
			respondWithError(w, http.StatusBadRequest, "bad request")
			log.Printf("bad request, error: %s,", err)

			return
		}

		beginDate, endDate := getRange(dates)

		log.Printf("range is from: %s to: %s", beginDate, endDate)

		start := time.Now().UTC()

		reviews, err := mongo.FindRange(beginDate, endDate)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			log.Printf("failed to get data from DB, error: %s,", err.Error())
			return
		}
		timeFetch := time.Since(start)
		start = time.Now().UTC()
		if len(reviews) == 0 {
			respondWithJSON(w, http.StatusOK, []StrShingleFreq{})
			log.Println("zero size list")
			return
		}

		// make chan of strings of shingles
		shingles := make(chan string)
		var postList sync.WaitGroup
		postList.Add(len(reviews))
		for i := range reviews {
			go func(review *ReviewIndex, i int) {
				defer postList.Done()
				post := ReviewPost{i, review.Post}
				sList := getShinglesList(post)
				for _, s := range sList {
					shingles <- s
				}
			}(&reviews[i], i)
		}
		// push shingles to list
		ls := make([]string, 0)
		go func() {
			for shingle := range shingles {
				ls = append(ls, shingle)
			}
		}()
		postList.Wait()
		timeList := time.Since(start)
		start = time.Now().UTC()

		// sort shingles by frequency
		sortedByFreq := sortPhrasesByFreq(ls)
		timeSort := time.Since(start)
		start = time.Now().UTC()

		// get top of list
		topShingles := make([]StrShingleFreq, topList)
		tail := len(sortedByFreq) - 1
		if tail < 0 {
			respondWithJSON(w, http.StatusOK, []StrShingleFreq{})
			log.Println("zero size list")
			return
		}
		for ix := 0; ix < topList; ix++ {
			topShingles[ix] = sortedByFreq[tail-ix]
			log.Printf("%v\n", topShingles[ix].String())
		}
		timePrint := time.Since(start)

		log.Printf("length of list of shingles: %d", len(sortedByFreq))
		log.Printf("done fetching data from mgo: %v", timeFetch)
		log.Printf("done making shingles list: %v", timeList)
		log.Printf("done sorting shingles list: %v", timeSort)
		log.Printf("done printing sorted list: %v", timePrint)

		respondWithJSON(w, http.StatusOK, topShingles)
	}
}

func getRange(dates dateRangePost) (time.Time, time.Time) {
	beginDate := time.Date(
		dates.BeginDate.Year,
		dates.BeginDate.Month,
		dates.BeginDate.Day, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(
		dates.EndDate.Year,
		dates.EndDate.Month,
		dates.EndDate.Day, 0, 0, 0, 0, time.UTC)

	return beginDate, endDate
}
