package main

import (
	"sync"
)

type sevenBallDrawMap map[int64]*SevenBallDraw

func NewSevenBallDrawResults() *SevenBallDrawResults {
	return &SevenBallDrawResults{
		results: make(map[string]sevenBallDrawMap),
	}
}

type SevenBallDrawResults struct {
	lock    sync.Mutex
	results map[string]sevenBallDrawMap
	total   int64
}

func (s *SevenBallDrawResults) Total() int64 {
	return s.total
}

func (s *SevenBallDrawResults) Add(date string, draw SevenBallDraw) {
	s.lock.Lock()

	drawList, ok := s.results[date]
	if !ok {
		drawList = make(sevenBallDrawMap)
		s.results[date] = drawList
	}

	_, ok = drawList[draw.ID]
	if !ok {
		// I'm adding a draw that did not exist before
		s.total++
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
