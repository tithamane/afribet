package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

type sevenBallDrawMap map[int64]*SevenBallDraw

func NewSevenBallDrawResults() *SevenBallDrawResults {
	return &SevenBallDrawResults{
		Results: make(map[string]sevenBallDrawMap),
	}
}

type SevenBallDrawResults struct {
	lock    sync.Mutex
	Results map[string]sevenBallDrawMap `json="results"`
	Total   int64                       `json="total"`
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
		s.Add(date, draw)
	}
}

type SevenBallDraw struct {
	ID       int64 `json:"id"`
	UnixTime int64 `json:"unix_time"`
	Numbers  []int `json:"numbers"`
}
