package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// Team represents a baseball team data.
type Team struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// MatchData represents baseball match data at a single point in time.
type MatchData struct {
	GameID      int  `json:"game_id"`
	TeamHome    Team `json:"team_home"`
	TeamAway    Team `json:"team_away"`
	ScoreHome   int  `json:"home_score"`
	ScoreAway   int  `json:"away_score"`
	TopInning   bool `json:"top_inning"`
	Out         int  `json:"out"`
	FirstBase   bool `json:"1st_base"`
	SecondBase  bool `json:"2nd_base"`
	ThirdBase   bool `json:"3rd_base"`
	InningCount int  `json:"inning_count"`
	PitcherID   int  `json:"pitcher_id"`
	BatterID    int  `json:"batter_id"`
}

var (
	currentMatch *MatchData
)

func readJSON(which int) (*MatchData, error) {
	path := fmt.Sprint("data/", which, "/data.json")

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.New(fmt.Sprint("error reading a json file:", err.Error()))
	}
	matchData := MatchData{}
	err = json.Unmarshal(data, &matchData)
	if err != nil {
		return nil, err
	}

	return &matchData, nil
}

func updateJSON(finished chan struct{}) {
	i := 1
	for {
		fmt.Println("5 seconds!")
		var err error
		currentMatch, err = readJSON(i)
		if err != nil {
			break
		}

		i++
		time.Sleep(5 * time.Second)
	}

	close(finished)
}

func main() {
	finished := make(chan struct{})

	go updateJSON(finished)

	http.HandleFunc("/data", handle)
	http.ListenAndServe(":3000", nil)

	<-finished

	fmt.Println("fileserver: served all files, shutdown")
}

func handle(writer http.ResponseWriter, req *http.Request) {
	js, err := json.Marshal(currentMatch)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.Write(js)
}
