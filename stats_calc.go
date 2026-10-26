package main

import (
	"fmt"
	"sort"
	"strings"
)

type TeamWinRate struct {
	Name          string
	NumberOfGames int
	WinRate       float32
}

type TeamWinStreak struct {
	Name   string
	Streak int
}

type SharedStats struct {
	Shortest      MatchData
	Longest       MatchData
	RadiantWin    float32
	DireWin       float32
	TopTeam       TeamWinRate
	LongestStreak TeamWinStreak
}

func mapToSlice(matchMap map[int]MatchData) []MatchData {
	matches := make([]MatchData, 0, len(matchMap))
	for _, match := range matchMap {
		matches = append(matches, match)
	}
	return matches
}

func GetAllTeamNames(matches []MatchData) map[string]struct{} {
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

func FindMatchesByTeam(matches []MatchData, teamName string) []MatchData {
	var results []MatchData
	for _, match := range matches {
		if strings.EqualFold(match.RadiantTeam, teamName) || strings.EqualFold(match.DireTeam, teamName) {
			results = append(results, match)
		}
	}
	return results
}

func FindMatchesByTournament(matches []MatchData, tournament string) []MatchData {
	var results []MatchData
	for _, match := range matches {
		if strings.EqualFold(match.LeagueName, tournament) {
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

func FindSideWinRate(matches []MatchData) (float32, float32) {
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

func FindHighestWinRateTeam(matches []MatchData) TeamWinRate {
	teams := GetAllTeamNames(matches)
	best := TeamWinRate{}

	for teamName := range teams {
		teamMatches := FindMatchesByTeam(matches, teamName)
		teamWinrate := FindTeamWinRate(teamMatches, teamName)

		if teamWinrate > best.WinRate || (teamWinrate == best.WinRate && len(teamMatches) > best.NumberOfGames) {
			best.Name = teamName
			best.NumberOfGames = len(teamMatches)
			best.WinRate = teamWinrate
		}
	}
	return best
}

func FindHighestWinStreak(matches []MatchData) TeamWinStreak {
	if len(matches) == 0 {
		return TeamWinStreak{}
	}

	sortedMatches := make([]MatchData, len(matches))
	copy(sortedMatches, matches)

	// Sort chronologically: oldest -> newest
	sort.Slice(sortedMatches, func(i, j int) bool {
		return sortedMatches[i].StartTime < sortedMatches[j].StartTime
	})

	currentStreaks := make(map[string]int)
	bestStreaks := make(map[string]int)

	for _, match := range sortedMatches {
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
			topStreak.Name = team
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

func formatSideWinRateComparison(radiantWin, direWin float32) string {
	if radiantWin > direWin {
		return fmt.Sprintf("Radiant had a higher winrate with: %.2f%%.", radiantWin)
	} else if radiantWin == direWin {
		return "Radiant and Dire both have a 50% winrate."
	}
	return fmt.Sprintf("Dire had a higher winrate with: %.2f%%.", direWin)
}

func formatDuration(totalSeconds int) string {
	minutes := totalSeconds / 60
	seconds := totalSeconds % 60
	return fmt.Sprintf("%d:%02d", minutes, seconds)
}

func computeSharedStats(matches []MatchData, includeStreak bool) SharedStats {
	if len(matches) == 0 {
		return SharedStats{}
	}

	shortestChan := make(chan MatchData, 1)
	longestChan := make(chan MatchData, 1)
	radiantChan := make(chan float32, 1)
	direChan := make(chan float32, 1)
	highestWinRateChan := make(chan TeamWinRate, 1)
	streakChan := make(chan TeamWinStreak, 1)

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
		highestWinRateChan <- FindHighestWinRateTeam(matches)
	}()

	if includeStreak {
		go func() {
			streakChan <- FindHighestWinStreak(matches)
		}()
	}

	stats := SharedStats{
		Shortest:   <-shortestChan,
		Longest:    <-longestChan,
		RadiantWin: <-radiantChan,
		DireWin:    <-direChan,
		TopTeam:    <-highestWinRateChan,
	}

	if includeStreak {
		stats.LongestStreak = <-streakChan
	}

	return stats
}
