package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type SevenBallProcessor struct {
	content   string
	parsed    bool
	Draws     []SevenBallDraw
	MorePages []string
}

func NewSevenBallProcessor(html string) SevenBallProcessor {
	return SevenBallProcessor{content: html, Draws: []SevenBallDraw{}}
}

func (s *SevenBallProcessor) GetContent() string {
	return s.content
}

func (s *SevenBallProcessor) Parse() error {
	if s.parsed == true {
		return nil
	}

	document, err := s.createHtmlDocument(s.content)
	if err != nil {
		return err
	}

	morePages := s.extractMorePageLinks(document)
	if len(morePages) < 8 {
		fmt.Println("Links less that 8\n------------")
		fmt.Println(s.content)
		fmt.Println("-----------------")
	}
	draws := s.extractDraws(document)
	s.Draws = append(s.Draws, draws...)
	s.MorePages = append(s.MorePages, morePages...)
	s.parsed = true

	// Eventually return nil
	return nil
}

func (s *SevenBallProcessor) createHtmlDocument(html string) (*goquery.Document, error) {
	htmlReader := strings.NewReader(html)
	document, err := goquery.NewDocumentFromReader(htmlReader)
	if err != nil {
		return nil, err
	}

	return document, nil
}

func (s *SevenBallProcessor) extractDraws(document *goquery.Document) []SevenBallDraw {
	draws := []SevenBallDraw{}
	rows := s.extractRows(document)
	date := s.extractDate(document)

	rows.Each(func(_ int, row *goquery.Selection) {
		rowTimeStr := s.getTime(row)
		timeStr := fmt.Sprintf("%sT%sZ", date, rowTimeStr)
		rowTime, err := time.Parse(time.RFC3339, timeStr)
		if err != nil {
			log.Printf("Count not convert \"%s\" into a Go time.Time\n", timeStr)
			return
		}

		unixTime := rowTime.Unix()
		rowNumber, err := strconv.Atoi(s.getID(row))
		if err != nil {
			log.Printf("Count not convert \"%s\" into a Go int64\n", timeStr)
			return
		}

		rowID := int64(rowNumber)
		numbers := s.getResults(row)

		draw := SevenBallDraw{
			ID:       rowID,
			UnixTime: unixTime,
			Numbers:  numbers,
		}

		draws = append(draws, draw)
	})

	return draws
}

func (s *SevenBallProcessor) extractDate(document *goquery.Document) string {
	period := document.Find("#period div").Text()
	meta := strings.Split(period, " ")
	date := meta[0]
	return date
}

func (s *SevenBallProcessor) extractRows(document *goquery.Document) *goquery.Selection {
	rows := document.Find(".table.table-results tbody tr")
	return rows
}

func (s *SevenBallProcessor) getTime(selection *goquery.Selection) string {
	// TODO: Subract 2 hours from the returned time
	firstTd := selection.Find("td").First().Text()
	tdData := strings.Split(firstTd, " - ")
	tdTime := tdData[0]

	return tdTime
}

func (s *SevenBallProcessor) getID(selection *goquery.Selection) string {
	firstTd := selection.Find("td").First().Text()
	tdData := strings.Split(firstTd, " - ")
	tdID := tdData[1]

	return tdID
}

func (s *SevenBallProcessor) getResults(selection *goquery.Selection) []int {
	spans := selection.Find("span span")
	numbers := []int{}
	spans.Each(func(_ int, s *goquery.Selection) {
		number, err := strconv.Atoi(s.Text())
		if err != nil {
			// TODO: Encode more information on which results it was extracting from
			log.Println("Error extract number from results")
			numbers = append(numbers, -1)
		} else {
			numbers = append(numbers, number)
		}

	})

	return numbers
}

func (s *SevenBallProcessor) extractMorePageLinks(document *goquery.Document) []string {
	anchors := document.Find(".pagination a")
	links := []string{}
	anchors.Each(func(_ int, s *goquery.Selection) {
		link, ok := s.Attr("href")
		if !ok {
			// TODO: Print more data about which page this error occured on
			log.Println("Could not extract link attribute from pagination anchor")
			return
		}

		links = append(links, link)
	})

	return links
}
