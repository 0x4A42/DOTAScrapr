package main

import (
	"fmt"
)

func genericStats(matchMap map[int]MatchData) {
	if len(matchMap) == 0 {
		fmt.Println("No match data available.")
		return
	}

	matches := mapToSlice(matchMap)

	shortestChan := make(chan MatchData, 1)
	longestChan := make(chan MatchData, 1)
	radiantChan := make(chan float32, 1)
	direChan := make(chan float32, 1)
	streakChan := make(chan TeamWinStreak, 1)
	highestWinRateChan := make(chan TeamWinRate, 1)

	go func() {
		shortestChan <- FindShortestMatch(matches)
	}()

	go func() {
		longestChan <- FindLongestMatch(matches)
	}()

	go func() {
		radiant, dire := FindSideWinRate(matches)
		radiantChan <- radiant
		direChan <- dire
	}()

	go func() {
		streakChan <- FindHighestWinStreak(matches)
	}()

	go func() {
		highestWinRateChan <- FindHighestWinRateTeam(matches)
	}()

	// Synchronize and receive concurrent metrics
	shortest := <-shortestChan
	longest := <-longestChan
	radiantWin := <-radiantChan
	direWin := <-direChan
	longestStreak := <-streakChan
	topTeam := <-highestWinRateChan

	printGenericStats(shortest, longest, radiantWin, direWin, longestStreak, topTeam)
}

func printGenericStats(shortest, longest MatchData, radiantWin, direWin float32, longestStreak TeamWinStreak, topTeam TeamWinRate) {
	fmt.Printf("======================================\n")
	fmt.Printf("Stats for last 100 pro matches\n")
	fmt.Printf("Shortest match was %s vs %s: %s minutes.\n", shortest.RadiantTeam, shortest.DireTeam, formatDuration(shortest.Duration))
	fmt.Printf("Longest match was %s vs %s: %s minutes.\n", longest.RadiantTeam, longest.DireTeam, formatDuration(longest.Duration))

	fmt.Println(formatSideWinRateComparison(radiantWin, direWin))

	if longestStreak.Streak > 0 {
		fmt.Printf("Longest win streak: %s with %d consecutive wins.\n", longestStreak.Team, longestStreak.Streak)
	}

	if topTeam.Team != "" {
		fmt.Printf("Highest win rate: %s with %.2f%% across %d games.\n", topTeam.Team, topTeam.WinRate, topTeam.NumberOfGames)
	}

	fmt.Printf("======================================\n")
}
