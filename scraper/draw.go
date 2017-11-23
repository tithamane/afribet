package main

type SevenBallDraw struct {
	ID       int64 `json:"id"`
	UnixTime int64 `json:"unix_time"`
	Numbers  []int `json:"numbers"`
}
