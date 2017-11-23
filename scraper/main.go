package main

func main() {
	s := NewSevenBallFetcher()
	_ = s.FetchByDate("2017-11-23")
}
