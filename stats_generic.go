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
	stats := computeSharedStats(matches, true)
	printGenericStats(stats)
}

func printGenericStats(stats SharedStats) {
	fmt.Printf("======================================\n")
	fmt.Printf("Stats for last 100 pro matches\n")
	fmt.Printf("Shortest match was %s vs %s: %s minutes.\n", stats.Shortest.RadiantTeam, stats.Shortest.DireTeam, formatDuration(stats.Shortest.Duration))
	fmt.Printf("Longest match was %s vs %s: %s minutes.\n", stats.Longest.RadiantTeam, stats.Longest.DireTeam, formatDuration(stats.Longest.Duration))

	fmt.Println(formatSideWinRateComparison(stats.RadiantWin, stats.DireWin))

	if stats.LongestStreak.Streak > 0 {
		fmt.Printf("Longest win streak: %s with %d consecutive wins.\n", stats.LongestStreak.Name, stats.LongestStreak.Streak)
	}

	if stats.TopTeam.Name != "" {
		fmt.Printf("Highest win rate: %s with %.2f%% across %d games.\n", stats.TopTeam.Name, stats.TopTeam.WinRate, stats.TopTeam.NumberOfGames)
	}

	fmt.Printf("======================================\n")
}
