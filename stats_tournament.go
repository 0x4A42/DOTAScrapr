package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func tournamentStats(matchMap map[int]MatchData) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("What tournament would you like stats for?")
	tournament, _ := reader.ReadString('\n')
	tournament = strings.TrimSpace(tournament)

	matches := mapToSlice(matchMap)
	tournamentMatches := FindMatchesByTournament(matches, tournament)

	if len(tournamentMatches) == 0 {
		fmt.Printf("No matches found for tournament: %s\n", tournament)
		return
	}

	shortestChan := make(chan MatchData, 1)
	longestChan := make(chan MatchData, 1)
	radiantChan := make(chan float32, 1)
	direChan := make(chan float32, 1)
	highestWinRateChan := make(chan TeamWinRate, 1)

	go func() {
		shortestChan <- FindShortestMatch(tournamentMatches)
	}()

	go func() {
		longestChan <- FindLongestMatch(tournamentMatches)
	}()

	go func() {
		radiant, dire := FindSideWinRate(tournamentMatches)
		radiantChan <- radiant
		direChan <- dire
	}()

	go func() {
		highestWinRateChan <- FindHighestWinRateTeam(tournamentMatches)
	}()

	// Synchronize and receive concurrent results
	shortest := <-shortestChan
	longest := <-longestChan
	radiantWin := <-radiantChan
	direWin := <-direChan
	topTeam := <-highestWinRateChan

	printTournamentStats(tournament, len(tournamentMatches), shortest, longest, radiantWin, direWin, topTeam)
}

func printTournamentStats(tournament string, totalMatches int, shortest, longest MatchData, radiantWin, direWin float32, topTeam TeamWinRate) {
	fmt.Printf("======================================\n")
	fmt.Printf("Stats for %s\n", tournament)
	fmt.Printf("Total matches: %d\n", totalMatches)
	fmt.Printf("Shortest match was %s vs %s: %s minutes.\n", shortest.RadiantTeam, shortest.DireTeam, formatDuration(shortest.Duration))
	fmt.Printf("Longest match was %s vs %s: %s minutes.\n", longest.RadiantTeam, longest.DireTeam, formatDuration(longest.Duration))

	fmt.Println(formatSideWinRateComparison(radiantWin, direWin))

	if topTeam.Team != "" {
		fmt.Printf("Highest win rate team: %s (%.2f%% across %d games)\n", topTeam.Team, topTeam.WinRate, topTeam.NumberOfGames)
	}

	fmt.Printf("======================================\n")
}
