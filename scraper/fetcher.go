package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// date - y-m-d
func formatAfricaGameUrl(date string, page int) string {
	return fmt.Sprintf("https://www.betgamesafrica.co.za/ext/game/results/testpartner/%s/%d", date, page)
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

	// Fetch results for the remaining pages
	for _, pageLink := range firstPageProcessor.MorePages {
		page := fetchPage(pageLink)
		processor := NewSevenBallProcessor(page)
		list = append(list, processor.Draws...)
	}

	return list
}

func fetchPage(url string) string {
	res, err := http.Get(url)
	if err != nil {
		log.Printf("Could not fetch contents of \"%s\"\n", url)
		return ""
	}

	contents, err := ioutil.ReadAll(res.Body)

	if err != nil {
		log.Printf("Could not read contents of \"%s\"\n", url)
		return ""
	}

	return string(contents[:])
}
