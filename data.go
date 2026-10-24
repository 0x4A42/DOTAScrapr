package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
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

type MatchCache struct {
	mu        sync.RWMutex
	data      map[int]MatchData
	expiresAt time.Time
	ttl       time.Duration
}

func NewMatchCache(ttl time.Duration) *MatchCache {
	return &MatchCache{
		ttl: ttl,
	}
}

func (c *MatchCache) GetData() (map[int]MatchData, error) {
	// check if valid cache exists
	c.mu.RLock()
	if c.data != nil && time.Now().Before(c.expiresAt) {
		defer c.mu.RUnlock()
		fmt.Println("Valid cache found")
		return c.data, nil
	}
	c.mu.RUnlock()

	// refresh cache
	c.mu.Lock()
	defer c.mu.Unlock()

	// Double-check condition
	if c.data != nil && time.Now().Before(c.expiresAt) {
		return c.data, nil
	}

	freshData, err := getMatchData()
	if err != nil {
		return nil, err
	}

	c.data = freshData
	c.expiresAt = time.Now().Add(c.ttl)

	return c.data, nil
}

func getMatchData() (map[int]MatchData, error) {
	resp, err := fetchMatchDataFromAPI()
	if err != nil {
		return nil, err
	}

	return convertAPIResponse(resp)
}
func fetchMatchDataFromAPI() (*http.Response, error) {
	resp, err := http.Get("https://api.opendota.com/api/proMatches")
	if err != nil {
		return nil, err
	}
	return resp, nil
}
func convertAPIResponse(resp *http.Response) (map[int]MatchData, error) {
	defer resp.Body.Close()

	var matches []MatchData
	if err := json.NewDecoder(resp.Body).Decode(&matches); err != nil {
		return nil, err
	}

	matchMap := make(map[int]MatchData, len(matches))
	for i, m := range matches {
		matchMap[i] = m
	}

	return matchMap, nil
}
