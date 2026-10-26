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

	stats := computeSharedStats(tournamentMatches, false)
	printTournamentStats(tournament, len(tournamentMatches), stats)
}

func printTournamentStats(tournament string, totalMatches int, stats SharedStats) {
	fmt.Printf("======================================\n")
	fmt.Printf("Stats for %s\n", tournament)
	fmt.Printf("Total matches: %d\n", totalMatches)
	fmt.Printf("Shortest match was %s vs %s: %s minutes.\n", stats.Shortest.RadiantTeam, stats.Shortest.DireTeam, formatDuration(stats.Shortest.Duration))
	fmt.Printf("Longest match was %s vs %s: %s minutes.\n", stats.Longest.RadiantTeam, stats.Longest.DireTeam, formatDuration(stats.Longest.Duration))

	fmt.Println(formatSideWinRateComparison(stats.RadiantWin, stats.DireWin))

	if stats.TopTeam.Name != "" {
		fmt.Printf("Highest win rate team: %s (%.2f%% across %d games)\n", stats.TopTeam.Name, stats.TopTeam.WinRate, stats.TopTeam.NumberOfGames)
	}

	fmt.Printf("======================================\n")
}
