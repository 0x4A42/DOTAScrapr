package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func teamStats(matchMap map[int]MatchData) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("What team would you like stats for? ")
	team, _ := reader.ReadString('\n')
	team = strings.TrimSpace(team)

	matches := mapToSlice(matchMap)
	teamMatches := FindMatchesByTeam(matches, team)

	if len(teamMatches) == 0 {
		fmt.Printf("%s did not play in the last 100 Pro games.\n", team)
		return
	}

	shortest := FindShortestMatch(teamMatches)
	longest := FindLongestMatch(teamMatches)
	winrate := FindTeamWinRate(teamMatches, team)
	currentStreak, maxStreak := FindTeamCurrentAndMaxStreak(teamMatches, team)

	printTeamStats(team, shortest, longest, winrate, len(teamMatches), currentStreak, maxStreak)
}

func printTeamStats(team string, shortest, longest MatchData, winrate float32, totalMatches, currentStreak, maxStreak int) {
	fmt.Printf("======================================\n")
	fmt.Printf("Stats for %s\n", team)
	fmt.Printf("Shortest match: %s minutes.\n", formatDuration(shortest.Duration))
	fmt.Printf("Longest match: %s minutes.\n", formatDuration(longest.Duration))
	fmt.Printf("Win rate: %.2f%% across %d matches\n", winrate, totalMatches)
	fmt.Printf("Current win streak: %d | Longest win streak: %d\n", currentStreak, maxStreak)
	fmt.Printf("======================================\n")
}
