package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

type TeamWinRate struct {
	Team          string
	NumberOfGames int
	WinRate       float32
}

type TeamWinStreak struct {
	Team   string
	Streak int
}

func teamStats(matchMap map[int]MatchData) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("What team would you like stats for? ")
	team, _ := reader.ReadString('\n')
	team = strings.TrimSpace(team)

	teamMatches := FindMatchesByTeam(matchMap, team)

	if len(teamMatches) == 0 {
		fmt.Printf("%s did not play in the last 100 Pro games.\n", team)
		return
	}

	shortest := FindShortestMatch(teamMatches)
	longest := FindLongestMatch(teamMatches)
	winrate := FindTeamWinRate(teamMatches, team)
	currentStreak, maxStreak := FindTeamCurrentAndMaxStreak(teamMatches, team)

	fmt.Printf("======================================\n")
	fmt.Printf("Stats for %s\n", team)
	fmt.Printf("Shortest match: %s minutes.\n", formatDuration(shortest.Duration))
	fmt.Printf("Longest match: %s minutes.\n", formatDuration(longest.Duration))
	fmt.Printf("Win rate: %.2f%% across %d matches\n", winrate, len(teamMatches))
	fmt.Printf("Current win streak: %d | Longest win streak: %d\n", currentStreak, maxStreak)
	fmt.Printf("======================================\n")
}

func tournamentStats(matchMap map[int]MatchData) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("What tournament would you like stats for?")
	tournament, _ := reader.ReadString('\n')
	tournament = strings.TrimSpace(tournament)
	tournamentMatches := FindMatchesByTournament(matchMap, tournament)

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
		radiant, dire := FindTournamentSideWinrate(tournamentMatches)
		radiantChan <- radiant
		direChan <- dire
	}()

	go func() {
		teams := GetAllTeamNamesTournament(tournamentMatches)
		best := TeamWinRate{}

		for team := range teams {
			tMatches := FindMatchesByTournamentAndTeam(tournamentMatches, team)
			rate := FindTeamWinRate(tMatches, team)

			if rate > best.WinRate || (rate == best.WinRate && len(tMatches) > best.NumberOfGames) {
				best.Team = team
				best.NumberOfGames = len(tMatches)
				best.WinRate = rate
			}
		}
		highestWinRateChan <- best
	}()

	// Synchronize and receive concurrent results
	shortest := <-shortestChan
	longest := <-longestChan
	radiantWin := <-radiantChan
	direWin := <-direChan
	topTeam := <-highestWinRateChan

	fmt.Printf("======================================\n")
	fmt.Printf("Stats for %s\n", tournament)
	fmt.Printf("Total matches: %d\n", len(tournamentMatches))
	fmt.Printf("Shortest match was %s vs %s: %s minutes.\n", shortest.RadiantTeam, shortest.DireTeam, formatDuration(shortest.Duration))
	fmt.Printf("Longest match was %s vs %s: %s minutes.\n", longest.RadiantTeam, longest.DireTeam, formatDuration(longest.Duration))

	if radiantWin > direWin {
		fmt.Printf("Radiant had a higher winrate with: %.2f%%.\n", radiantWin)
	} else if radiantWin == direWin {
		fmt.Println("Radiant and Dire both have a 50% winrate.")
	} else {
		fmt.Printf("Dire had a higher winrate with: %.2f%%.\n", direWin)
	}

	if topTeam.Team != "" {
		fmt.Printf("Highest win rate team: %s (%.2f%% across %d games)\n", topTeam.Team, topTeam.WinRate, topTeam.NumberOfGames)
	}

	fmt.Printf("======================================\n")
}

func genericStats(matchMap map[int]MatchData) {
	if len(matchMap) == 0 {
		fmt.Println("No match data available.")
		return
	}

	shortestChan := make(chan MatchData, 1)
	longestChan := make(chan MatchData, 1)
	radiantChan := make(chan float32, 1)
	direChan := make(chan float32, 1)
	streakChan := make(chan TeamWinStreak, 1)
	highestWinRateChan := make(chan TeamWinRate, 1)

	go func() {
		shortestChan <- FindShortestMatchByMap(matchMap)
	}()

	go func() {
		longestChan <- FindLongestMatchByMap(matchMap)
	}()

	go func() {
		radiant, dire := FindSideWinRate(matchMap)
		radiantChan <- radiant
		direChan <- dire
	}()

	go func() {
		streakChan <- FindHighestWinStreak(matchMap)
	}()

	go func() {
		teams := GetAllTeamNames(matchMap)
		best := TeamWinRate{}

		for teamName := range teams {
			teamMatches := FindMatchesByTeam(matchMap, teamName)
			teamWinrate := FindTeamWinRate(teamMatches, teamName)

			if teamWinrate > best.WinRate || (teamWinrate == best.WinRate && len(teamMatches) > best.NumberOfGames) {
				best.Team = teamName
				best.NumberOfGames = len(teamMatches)
				best.WinRate = teamWinrate
			}
		}
		highestWinRateChan <- best
	}()

	// Synchronize and receive concurrent metrics
	shortest := <-shortestChan
	longest := <-longestChan
	radiantWin := <-radiantChan
	direWin := <-direChan
	longestStreak := <-streakChan
	topTeam := <-highestWinRateChan

	fmt.Printf("======================================\n")
	fmt.Printf("Stats for last 100 pro matches\n")
	fmt.Printf("Shortest match was %s vs %s: %s minutes.\n", shortest.RadiantTeam, shortest.DireTeam, formatDuration(shortest.Duration))
	fmt.Printf("Longest match was %s vs %s: %s minutes.\n", longest.RadiantTeam, longest.DireTeam, formatDuration(longest.Duration))

	if radiantWin > direWin {
		fmt.Printf("Radiant had a higher winrate with: %.2f%%.\n", radiantWin)
	} else if radiantWin == direWin {
		fmt.Println("Radiant and Dire both have a 50% winrate.")
	} else {
		fmt.Printf("Dire had a higher winrate with: %.2f%%.\n", direWin)
	}

	if longestStreak.Streak > 0 {
		fmt.Printf("Longest win streak: %s with %d consecutive wins.\n", longestStreak.Team, longestStreak.Streak)
	}

	if topTeam.Team != "" {
		fmt.Printf("Highest win rate: %s with %.2f%% across %d games.\n", topTeam.Team, topTeam.WinRate, topTeam.NumberOfGames)
	}

	fmt.Printf("======================================\n")
}

func GetAllTeamNames(matchMap map[int]MatchData) map[string]struct{} {
	teams := make(map[string]struct{})
	for _, match := range matchMap {
		if match.RadiantTeam != "" {
			teams[match.RadiantTeam] = struct{}{}
		}
		if match.DireTeam != "" {
			teams[match.DireTeam] = struct{}{}
		}
	}
	return teams
}

func GetAllTeamNamesTournament(matches []MatchData) map[string]struct{} {
	teams := make(map[string]struct{})
	for _, match := range matches {
		if match.RadiantTeam != "" {
			teams[match.RadiantTeam] = struct{}{}
		}
		if match.DireTeam != "" {
			teams[match.DireTeam] = struct{}{}
		}
	}
	return teams
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

func FindMatchesByTournament(matchMap map[int]MatchData, tournament string) []MatchData {
	var results []MatchData
	for _, match := range matchMap {
		if strings.EqualFold(match.LeagueName, tournament) {
			results = append(results, match)
		}
	}
	return results
}

func FindMatchesByTournamentAndTeam(matches []MatchData, teamName string) []MatchData {
	var results []MatchData
	for _, match := range matches {
		if strings.EqualFold(match.RadiantTeam, teamName) || strings.EqualFold(match.DireTeam, teamName) {
			results = append(results, match)
		}
	}
	return results
}

func FindShortestMatch(matches []MatchData) MatchData {
	if len(matches) == 0 {
		return MatchData{}
	}
	shortest := matches[0]
	for _, match := range matches[1:] {
		if match.Duration < shortest.Duration {
			shortest = match
		}
	}
	return shortest
}

func FindShortestMatchByMap(matchMap map[int]MatchData) MatchData {
	if len(matchMap) == 0 {
		return MatchData{}
	}
	var shortest MatchData
	first := true
	for _, match := range matchMap {
		if first || match.Duration < shortest.Duration {
			shortest = match
			first = false
		}
	}
	return shortest
}

func FindLongestMatch(matches []MatchData) MatchData {
	if len(matches) == 0 {
		return MatchData{}
	}
	longest := matches[0]
	for _, match := range matches[1:] {
		if match.Duration > longest.Duration {
			longest = match
		}
	}
	return longest
}

func FindLongestMatchByMap(matchMap map[int]MatchData) MatchData {
	if len(matchMap) == 0 {
		return MatchData{}
	}
	var longest MatchData
	first := true
	for _, match := range matchMap {
		if first || match.Duration > longest.Duration {
			longest = match
			first = false
		}
	}
	return longest
}

func FindTeamWinRate(matches []MatchData, team string) float32 {
	if len(matches) == 0 {
		return 0.0
	}

	teamWins := 0
	for _, match := range matches {
		isRadiant := strings.EqualFold(match.RadiantTeam, team)
		isDire := strings.EqualFold(match.DireTeam, team)

		if (isRadiant && match.RadiantWin) || (isDire && !match.RadiantWin) {
			teamWins++
		}
	}

	return (float32(teamWins) / float32(len(matches))) * 100.0
}

func FindSideWinRate(matchMap map[int]MatchData) (float32, float32) {
	if len(matchMap) == 0 {
		return 0, 0
	}

	radiantWins := 0
	direWins := 0

	for _, match := range matchMap {
		if match.RadiantWin {
			radiantWins++
		} else {
			direWins++
		}
	}

	radiantRate := (float32(radiantWins) / float32(len(matchMap))) * 100.0
	direRate := (float32(direWins) / float32(len(matchMap))) * 100.0

	return radiantRate, direRate
}

func FindTournamentSideWinrate(matches []MatchData) (float32, float32) {
	if len(matches) == 0 {
		return 0, 0
	}

	radiantWins := 0
	direWins := 0

	for _, match := range matches {
		if match.RadiantWin {
			radiantWins++
		} else {
			direWins++
		}
	}

	radiantRate := (float32(radiantWins) / float32(len(matches))) * 100.0
	direRate := (float32(direWins) / float32(len(matches))) * 100.0

	return radiantRate, direRate
}

func FindHighestWinStreak(matchMap map[int]MatchData) TeamWinStreak {
	if len(matchMap) == 0 {
		return TeamWinStreak{}
	}

	matches := make([]MatchData, 0, len(matchMap))
	for _, match := range matchMap {
		matches = append(matches, match)
	}

	// Sort chronologically: oldest -> newest
	sort.Slice(matches, func(i, j int) bool {
		return matches[i].StartTime < matches[j].StartTime
	})

	currentStreaks := make(map[string]int)
	bestStreaks := make(map[string]int)

	for _, match := range matches {
		var winner, loser string

		if match.RadiantWin {
			winner = match.RadiantTeam
			loser = match.DireTeam
		} else {
			winner = match.DireTeam
			loser = match.RadiantTeam
		}

		if winner != "" {
			currentStreaks[winner]++
			if currentStreaks[winner] > bestStreaks[winner] {
				bestStreaks[winner] = currentStreaks[winner]
			}
		}

		if loser != "" {
			currentStreaks[loser] = 0
		}
	}

	var topStreak TeamWinStreak
	for team, streak := range bestStreaks {
		if streak > topStreak.Streak {
			topStreak.Team = team
			topStreak.Streak = streak
		}
	}

	return topStreak
}

func FindTeamCurrentAndMaxStreak(teamMatches []MatchData, teamName string) (int, int) {
	if len(teamMatches) == 0 {
		return 0, 0
	}

	sortedMatches := make([]MatchData, len(teamMatches))
	copy(sortedMatches, teamMatches)

	// Sort team matches chronologically: oldest -> newest
	sort.Slice(sortedMatches, func(i, j int) bool {
		return sortedMatches[i].StartTime < sortedMatches[j].StartTime
	})

	currentStreak := 0
	maxStreak := 0

	for _, match := range sortedMatches {
		isRadiant := strings.EqualFold(match.RadiantTeam, teamName)
		isDire := strings.EqualFold(match.DireTeam, teamName)

		won := (isRadiant && match.RadiantWin) || (isDire && !match.RadiantWin)

		if won {
			currentStreak++
			if currentStreak > maxStreak {
				maxStreak = currentStreak
			}
		} else {
			currentStreak = 0
		}
	}

	return currentStreak, maxStreak
}

func formatDuration(totalSeconds int) string {
	minutes := totalSeconds / 60
	seconds := totalSeconds % 60
	return fmt.Sprintf("%d:%02d", minutes, seconds)
}
