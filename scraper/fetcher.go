package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	SevenBallGame = 1
)

// date - y-m-d
func formatAfricaGameUrl(date string, page int) string {

	return fmt.Sprintf("https://www.betgamesafrica.co.za/ext/game/results/testpartner/%s/%d/%d", date, SevenBallGame, page)
}

func NewSevenBallFetcher() SevenBallFetcher {
	return SevenBallFetcher{}
}

type SevenBallFetcher struct {
}

func (s *SevenBallFetcher) FetchByDate(date string) []SevenBallDraw {
	// TODO: Check for errors if processing fails
	list := []SevenBallDraw{}
	sevenBallUrl := formatAfricaGameUrl(date, 1)
	firstPage := fetchPage(sevenBallUrl)
	firstPageProcessor := NewSevenBallProcessor(firstPage)
	_ = firstPageProcessor.Parse()
	list = append(list, firstPageProcessor.Draws...)

	// Download 7 other pages if date is not today,
	// the html for some dates is inconsistent and needs investigating
	now := time.Now()
	today := fmt.Sprintf("%d-%02d-%02d", now.Year(), now.Month(), now.Day())
	var pageLinks []string

	if strings.Compare(date, today) != 0 {
		for i := 2; i <= 8; i++ {
			link := formatAfricaGameUrl(date, i)
			pageLinks = append(pageLinks, link)
		}
	} else {
		pageLinks = firstPageProcessor.MorePages
	}

	// Fetch results for the remaining pages
	for _, pageLink := range pageLinks {
		page := fetchPage(pageLink)
		processor := NewSevenBallProcessor(page)
		_ = processor.Parse()
		list = append(list, processor.Draws...)
	}

	fmt.Println("Date:", date, " NumberOfResults:", len(list))

	return list
}

func fetchPage(url string) string {
	res, err := http.Get(url)
	if err != nil {
		log.Printf("Could not fetch contents of \"%s\"\n\t>> %s\n", url, err)
		return ""
	}

	contents, err := ioutil.ReadAll(res.Body)

	if err != nil {
		log.Printf("Could not read contents of \"%s\"\n", url)
		return ""
	}

	return string(contents[:])
}
