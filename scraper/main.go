package main

import (
	"fmt"
)

func main() {
	s := NewSevenBallFetcher()
	_ = s.FetchByDate("2017-11-23")

	dr, err := NewDateRange("2016-02-26", "2016-03-03")
	if err != nil {
		fmt.Println("You provided a faulty data:", err)
	} else {
		fmt.Println(dr.DatesBetween())
	}
}
