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

	teamMatches := FindMatchesByTeam(matchMap, team)

	if len(teamMatches) == 0 {
		fmt.Printf("%s did not play in the last 100 Pro games. \n", team)
		return
	}

	shortest := FindShortestMatch(teamMatches)
	longest := FindLongestMatch(teamMatches)
	winrate := FindTeamWinRate(teamMatches, team)

	fmt.Printf("======================================\n")
	fmt.Printf("Stats for %s \n", team)
	fmt.Printf("Shortest match: %d minutes. \n", shortest.Duration/60)
	fmt.Printf("Longest match: %d minutes. \n", longest.Duration/60)
	fmt.Printf("Win rate: %.2f%%\n across %d matches", winrate, len(teamMatches))
}

func tournamentStats(matchMap map[int]MatchData) {
	fmt.Println("What tournament would you like stats for?")
	var input string
	fmt.Scanln(&input)
}

func genericTeamStats(matchMap map[int]MatchData) {
	// most wins
	// longest game
	// shortest game
	// win streak?
}

func genericGameStats(matchMap map[int]MatchData) {
	// longest game
	// shortest game
	// radiant/dire win rate
}

func FindMatchesByTeam(matchMap map[int]MatchData, teamName string) []MatchData {
	var results []MatchData

	for _, match := range matchMap {
		if strings.EqualFold(match.RadiantTeam, teamName) || strings.EqualFold(match.DireTeam, teamName) {
			results = append(results, match)
		}
	}
	return results
}

func FindShortestMatch(matches []MatchData) MatchData {
	shortest := matches[0]
	for _, match := range matches {
		if match.Duration < shortest.Duration {
			shortest = match
		}
	}
	return shortest
}

func FindLongestMatch(matches []MatchData) MatchData {
	longest := matches[0]
	for _, match := range matches {
		if match.Duration > longest.Duration {
			longest = match
		}
	}
	return longest
}

func FindTeamWinRate(matches []MatchData, team string) float32 {
	teamWins := 0

	for _, match := range matches {
		if strings.EqualFold(match.RadiantTeam, team) {
			if match.RadiantWin {
				teamWins++
			}
		}
	}
	return float32((teamWins / len(matches)) * 100)
}

func FindSideWinRate(matchMap map[int]MatchData) (int, int) {
	radiantWins := 0
	direWins := 0

	for _, match := range matchMap {
		switch match.RadiantWin {
		case true:
			radiantWins++
		case false:
			direWins++
		}
	}

	return radiantWins, direWins
}
