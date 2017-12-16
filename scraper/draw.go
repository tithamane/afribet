package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
)

func NthNumber(n int) string {
	var title string
	switch n {
	case 1:
		title = "first"
	case 2:
		title = "second"
	case 3:
		title = "third"
	case 4:
		title = "fourth"
	case 5:
		title = "fifth"
	case 6:
		title = "sixth"
	case 7:
		title = "seventh"
	default:
		title = "unknown"
	}

	return fmt.Sprintf("%s_number", title)
}

type sevenBallDrawMap map[int64]*SevenBallDraw

func NewSevenBallDrawResults() *SevenBallDrawResults {
	return &SevenBallDrawResults{
		Results: make(map[string]sevenBallDrawMap),
	}
}

type SevenBallDrawResults struct {
	lock    sync.Mutex
	Results map[string]sevenBallDrawMap `json:"results"`
	Total   int64                       `json:"total"`
}

func (s *SevenBallDrawResults) SaveJSON() error {
	dataFolder := "data/scraped"
	err := os.MkdirAll(dataFolder, os.ModePerm)
	if os.IsNotExist(err) {
		log.Println(err)
		return err
	}

	data, err := json.Marshal(s)
	if err != nil {
		log.Println(err)
		return err
	}

	filename := fmt.Sprintf("%s/seven_ball_results.json", dataFolder)
	err = ioutil.WriteFile(filename, data, 0777)

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (s *SevenBallDrawResults) SaveCSV() error {
	dataPath := "data/csv"

	// Make the os directories if they don't exist
	err := os.MkdirAll(dataPath, os.ModePerm)
	if os.IsNotExist(err) {
		log.Println(err)
		return err
	}

	csvFilename := fmt.Sprintf("%s/seven_ball.csv", dataPath)
	csvFile, err := os.OpenFile(csvFilename, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Println(err)
		return err
	}
	defer csvFile.Close()

	titles := "draw_id,unix_time"
	numbersAsStr := []string{}
	// TODO: Change the hardcoded 7 to something more dynamic
	for i := 0; i < 7; i++ {
		number := NthNumber(i + 1)
		numbersAsStr = append(numbersAsStr, number)
	}
	numbers := strings.Join(numbersAsStr, ",")
	titles = fmt.Sprintf("%s,%s", titles, numbers)
	titleLine := []byte(fmt.Sprintln(titles))
	csvFile.Write(titleLine)

	for _, dayResults := range s.Results {
		for _, draw := range dayResults {
			line := fmt.Sprintln(draw.ToCSVLine())
			csvLine := []byte(line)
			csvFile.Write(csvLine)
		}
	}

	return nil
}

func (s *SevenBallDrawResults) SaveCombinationsCSV() {
	nCombinations := []int{2, 3, 4}
	var wg sync.WaitGroup
	wg.Add(len(nCombinations))

	for _, nCombination := range nCombinations {
		go func(n int) {
			s.SaveNCombinationCSV(n)
			wg.Done()
		}(nCombination)
	}
	wg.Wait()
	log.Printf("Done: NCombinations have been saved for %+v.\n", nCombinations)
}

func (s *SevenBallDrawResults) SaveNCombinationCSV(n int) {
	dataPath := "data/csv/combinations"

	// Make the os directories if they don't exist
	err := os.MkdirAll(dataPath, os.ModePerm)
	if os.IsNotExist(err) {
		log.Println(err)
		return
	}

	csvFilename := fmt.Sprintf("%s/%d_seven_ball.csv", dataPath, n)
	csvFile, err := os.OpenFile(csvFilename, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Println(err)
		return
	}
	defer csvFile.Close()

	// Save the results into the created file
	titles := "draw_id,length,combination_index,unix_time"
	numbersAsStr := []string{}
	// TODO: Create labels for the numbers that were part of the results
	for i := 0; i < n; i++ {
		number := NthNumber(i + 1)
		numbersAsStr = append(numbersAsStr, number)
	}
	numbers := strings.Join(numbersAsStr, ",")
	titles = fmt.Sprintf("%s,%s", titles, numbers)
	titleLine := []byte(fmt.Sprintln(titles))
	csvFile.Write(titleLine)

	for _, dayResults := range s.Results {
		for _, draw := range dayResults {
			combinations := draw.CombinationsN(n)
			for _, combination := range combinations {
				line := fmt.Sprintln(combination.ToCSVLine())
				csvLine := []byte(line)
				csvFile.Write(csvLine)
			}
		}
	}
}

func (s *SevenBallDrawResults) Add(date string, draw SevenBallDraw) {
	s.lock.Lock()

	drawList, ok := s.Results[date]
	if !ok {
		drawList = make(sevenBallDrawMap)
		s.Results[date] = drawList
	}

	_, ok = drawList[draw.ID]
	if !ok {
		// I'm adding a draw that did not exist before
		s.Total++
	}

	drawList[draw.ID] = &draw

	s.lock.Unlock()
}

func (s *SevenBallDrawResults) AddList(date string, draws []SevenBallDraw) {
	for _, draw := range draws {
		// TODO: Handle error where a draw has no numbers
		drawNumbersLength := len(draw.Numbers)
		if drawNumbersLength > 0 {
			s.Add(date, draw)
		}
	}
}

type SevenBallDraw struct {
	ID       int64 `json:"id"`
	UnixTime int64 `json:"unix_time"`
	Numbers  []int `json:"numbers"`
}

func (s *SevenBallDraw) ToCSVLine() string {
	meta := fmt.Sprintf("%d,%d", s.ID, s.UnixTime)

	numbersAsStr := []string{}
	for _, value := range s.Numbers {
		number := fmt.Sprintf("%d", value)
		numbersAsStr = append(numbersAsStr, number)
	}

	numbers := strings.Join(numbersAsStr, ",")
	line := fmt.Sprintf("%s,%s", meta, numbers)
	return line
}

func (s *SevenBallDraw) CombinationsN(n int) []DrawCombination {
	bitShifter, err := NewBitShifter(s.Numbers)
	if err != nil {
		log.Printf("Could create BitShifter for SevenBallDraw: %+v", *s)
		return []DrawCombination{}
	}

	combinations, err := bitShifter.CombinationsN(n)
	if err != nil {
		log.Printf("Could not generate CombinationsN(%d) for %+v\n", n, *s)
		return []DrawCombination{}
	}

	drawCombinations := []DrawCombination{}
	for i, combination := range combinations {
		drawCombination := DrawCombination{
			ID:               s.ID,
			Length:           n,
			CombinationIndex: i,
			UnixTime:         s.UnixTime,
			Numbers:          combination,
		}

		drawCombinations = append(drawCombinations, drawCombination)
	}

	return drawCombinations
}

type DrawCombination struct {
	ID               int64 `json:"draw_id"`
	Length           int   `json:"length"`
	CombinationIndex int   `json:"combination_index"`
	UnixTime         int64 `json:"unix_time"`
	Numbers          []int `json:"numbers"`
}

func (d *DrawCombination) CSVTitles(n int) string {
	meta := "draw_id,length,combination_index,unix_time"

	numbersAsStr := []string{}
	for i := 0; i < n; i++ {
		number := NthNumber(i + 1)
		numbersAsStr = append(numbersAsStr, number)
	}
	numberTitles := strings.Join(numbersAsStr, ",")
	titles := fmt.Sprintf("%s,%s", meta, numberTitles)
	return titles
}

func (d *DrawCombination) ToCSVLine() string {
	meta := fmt.Sprintf("%d,%d,%d,%d", d.ID, d.Length, d.CombinationIndex, d.UnixTime)

	numbersAsStr := []string{}
	for _, number := range d.Numbers {
		numbersAsStr = append(numbersAsStr, fmt.Sprintf("%d", number))
	}
	numbers := strings.Join(numbersAsStr, ",")

	finalStr := fmt.Sprintf("%s,%s", meta, numbers)
	return finalStr
}
