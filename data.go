package main

import (
	"encoding/json"
	"net/http"
)

type MatchData struct {
	Id           int    `json:"match_id"`
	Duration     int    `json:"duration"`
	StartTime    int    `json:"start_time"`
	RadiantTeam  string `json:"radiant_name"`
	DireTeam     string `json:"dire_name"`
	LeagueName   string `json:"league_name"`
	SeriesId     int    `json:"series_id"`
	RadiantScore int    `json:"radiant_score"`
	DireScore    int    `json:"dire_score"`
	RadiantWin   bool   `json:"radiant_win"`
}

func getMatchData() map[int]MatchData {
	client := http.Client{}
	resp, err := client.Get("https://api.opendota.com/api/proMatches")

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	var matches []MatchData
	if err := json.NewDecoder(resp.Body).Decode(&matches); err != nil {
		panic(err)
	}

	matchMap := make(map[int]MatchData, len(matches))
	for i, m := range matches {
		matchMap[i] = m
	}

	return matchMap
}
