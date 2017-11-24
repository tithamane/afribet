package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

const (
	Port               = 9292
	MaxDayDownloadable = 3700
)

var sevenBallDrawResults *SevenBallDrawResults

func main() {

	log.Println("Booting scraper")
	sevenBallDrawResults = NewSevenBallDrawResults()

	// s := NewSevenBallFetcher()
	// _ = s.FetchByDate("2017-11-23")

	// dr, err := NewDateRange("2016-02-26", "2016-03-03")
	// if err != nil {
	// 	fmt.Println("You provided a faulty data:", err)
	// } else {
	// 	fmt.Println(dr.DatesBetween())
	// }

	router := mux.NewRouter()
	setupRestAPI(router)

	port := fmt.Sprintf(":%d", Port)
	http.ListenAndServe(port, router)
}

func setupRestAPI(router *mux.Router) {
	router.HandleFunc("/fetch-between", FetchResultsBetween).Methods("POST")
	router.HandleFunc("/save-seven-ball-results", SaveSevenBallDrawResults)
}

func FetchResultsBetween(w http.ResponseWriter, r *http.Request) {
	// TODO: Don't fetch results if date is not today and they've been fetched already
	// TODO: Return errors as a object of some kind
	params := make(map[string]string)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintln(w, err)
		fmt.Println(err)
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal(body, &params)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintln(w, err)
		fmt.Println(err)
		return
	}

	from, ok := params["from"]
	if !ok {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintln(w, "Missing 'from' attribute")
		fmt.Println(err)
		return
	}

	to, ok := params["to"]
	if !ok {
		now := time.Now()
		to = fmt.Sprintf("%d-%02d-%02d", now.Year(), now.Month(), now.Day())
	}

	dateRange, err := NewDateRange(from, to)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintln(w, "Invalid date was used")
		fmt.Println(err)
		return
	}

	datesBetween := dateRange.DatesBetween()

	var wg sync.WaitGroup
	wg.Add(len(datesBetween))
	sevenBallFetcher := NewSevenBallFetcher()
	log.Println("Starting to process draws. DatesBetween:", len(datesBetween))
	for _, date := range datesBetween {
		// TODO: Check if the date is not today and has been downloaded, if so, ignore it
		go func(downloadDate string) {
			draws := sevenBallFetcher.FetchByDate(downloadDate)
			sevenBallDrawResults.AddList(downloadDate, draws)
			wg.Done()
		}(date)
	}
	wg.Wait()
	log.Println("Finished processing draws")
	fmt.Fprintf(w, "Done!")
}

func SaveSevenBallDrawResults(w http.ResponseWriter, r *http.Request) {
	log.Println("Saving seven ball results")
	err := sevenBallDrawResults.SaveJSON()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
	}

	w.WriteHeader(http.StatusOK)
}
