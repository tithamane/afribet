package main

import (
	"fmt"
	"time"
)

func NewDateRange(from, to string) (DateRange, error) {
	firstDateStr := fmt.Sprintf("%sT00:00:00Z", from)
	secondDateStr := fmt.Sprintf("%sT00:00:00Z", to)

	firstDate, err := time.Parse(time.RFC3339, firstDateStr)
	if err != nil {
		return DateRange{}, err
	}
	secondDate, err := time.Parse(time.RFC3339, secondDateStr)
	if err != nil {
		return DateRange{}, err
	}

	var toDate time.Time
	var fromDate time.Time
	if firstDate.After(secondDate) {
		toDate = firstDate
		fromDate = secondDate
	} else {
		toDate = secondDate
		fromDate = firstDate
	}

	return DateRange{from: fromDate, to: toDate}, nil
}

type DateRange struct {
	from time.Time
	to   time.Time
}

func (d *DateRange) DatesBetween() []string {
	dates := []string{}
	iteratorDate := d.from
	for {
		date := fmt.Sprintf("%d-%02d-%02d", iteratorDate.Year(), iteratorDate.Month(), iteratorDate.Day())
		dates = append(dates, date)
		iteratorDate = iteratorDate.Add(time.Hour * 24)
		if iteratorDate.After(d.to) {
			break
		}
	}
	return dates
}
