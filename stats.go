package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type TeamWinRate struct {
	Team          string
	NumberOfGames int
	WinRate       float32
}

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
	fmt.Printf("Win rate: %.2f%% across %d matches \n", winrate, len(teamMatches))
	fmt.Printf("======================================\n")
}

func tournamentStats(matchMap map[int]MatchData) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("What tournament would you like stats for?")
	tournament, _ := reader.ReadString('\n')
	tournament = strings.TrimSpace(tournament)
	tournamentMatches := FindMatchesByTournament(matchMap, tournament)

	shortestChan := make(chan MatchData, 1)
	longestChan := make(chan MatchData, 1)
	radiantChan := make(chan float32, 1)
	direChan := make(chan float32, 1)
	teamsChan := make(chan map[string]struct{}, 1)
	highestWinRate := TeamWinRate{}

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
		teamsChan <- GetAllTeamNamesTournament(tournamentMatches)
		teams := <-teamsChan

		for team := range teams {
			teamMatches := FindMatchesByTeam(matchMap, team)
			teamWinrate := FindTeamWinRate(teamMatches, team)

			if teamWinrate > highestWinRate.WinRate {
				highestWinRate.Team = team
				highestWinRate.NumberOfGames = len(teamMatches)
				highestWinRate.WinRate = teamWinrate
			}
		}
	}()

	// Receive results
	shortest := <-shortestChan
	longest := <-longestChan
	radiantWin := <-radiantChan
	direWin := <-direChan

	// win streak?

	fmt.Printf("Shortest match was %s vs %s: %s minutes.\n", shortest.RadiantTeam, shortest.DireTeam, formatDuration(shortest.Duration))
	fmt.Printf("Longest match was %s vs %s: %s minutes.\n", longest.RadiantTeam, longest.DireTeam, formatDuration(longest.Duration))

	if radiantWin > direWin {
		fmt.Printf("Radiant had a higher winrate with: %.2f%%.\n", radiantWin)
	} else if radiantWin == direWin {
		fmt.Println("Radiant and Dire both have a 50% winrate.")
	} else {
		fmt.Printf("Dire had a higher winrate with: %.2f%%.\n", direWin)
	}
}

func genericStats(matchMap map[int]MatchData) {

	shortestChan := make(chan MatchData, 1)
	longestChan := make(chan MatchData, 1)
	radiantChan := make(chan float32, 1)
	direChan := make(chan float32, 1)
	teamsChan := make(chan map[string]struct{}, 1)
	highestWinRate := TeamWinRate{}

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
		teamsChan <- GetAllTeamNames(matchMap)
		teams := <-teamsChan

		for teamName := range teams {
			teamMatches := FindMatchesByTeam(matchMap, teamName)
			teamWinrate := FindTeamWinRate(teamMatches, teamName)

			if teamWinrate > highestWinRate.WinRate {
				highestWinRate.Team = teamName
				highestWinRate.NumberOfGames = len(teamMatches)
				highestWinRate.WinRate = teamWinrate
			}
		}
	}()

	// Receive results
	shortest := <-shortestChan
	longest := <-longestChan
	radiantWin := <-radiantChan
	direWin := <-direChan

	// win streak?

	fmt.Printf("Shortest match was %s vs %s: %s minutes.\n", shortest.RadiantTeam, shortest.DireTeam, formatDuration(shortest.Duration))
	fmt.Printf("Longest match was %s vs %s: %s minutes.\n", longest.RadiantTeam, longest.DireTeam, formatDuration(longest.Duration))

	if radiantWin > direWin {
		fmt.Printf("Radiant had a higher winrate with: %.2f%%.\n", radiantWin)
	} else if radiantWin == direWin {
		fmt.Println("Radiant and Dire both have a 50% winrate.")
	} else {
		fmt.Printf("Dire had a higher winrate with: %.2f%%.\n", direWin)
	}

}

func GetAllTeamNames(matchMap map[int]MatchData) map[string]struct{} {

	teams := make(map[string]struct{})
	for _, match := range matchMap {
		_, hasRadiant := teams[match.RadiantTeam]
		_, hasDire := teams[match.DireTeam]

		if !hasRadiant {
			teams[match.RadiantTeam] = struct{}{}
		}
		if !hasDire {
			teams[match.DireTeam] = struct{}{}
		}
	}

	return teams
}

func GetAllTeamNamesTournament(matches []MatchData) map[string]struct{} {

	teams := make(map[string]struct{})
	for _, match := range matches {
		_, hasRadiant := teams[match.RadiantTeam]
		_, hasDire := teams[match.DireTeam]

		if !hasRadiant {
			teams[match.RadiantTeam] = struct{}{}
		}
		if !hasDire {
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

func FindShortestMatch(matches []MatchData) MatchData {
	shortest := matches[0]
	for _, match := range matches {
		if match.Duration < shortest.Duration {
			shortest = match
		}
	}
	return shortest
}

func FindShortestMatchByMap(matchMap map[int]MatchData) MatchData {
	shortest := matchMap[0]
	for _, match := range matchMap {
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

func FindLongestMatchByMap(matchMap map[int]MatchData) MatchData {
	longest := matchMap[0]
	for _, match := range matchMap {
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

func FindSideWinRate(matchMap map[int]MatchData) (float32, float32) {
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

	radiantRate := (float32(radiantWins) / float32(len(matchMap))) * 100
	direRate := (float32(direWins) / float32(len(matchMap))) * 100

	return radiantRate, direRate
}

func FindTournamentSideWinrate(matches []MatchData) (float32, float32) {
	radiantWins := 0
	direWins := 0

	for _, match := range matches {
		switch match.RadiantWin {
		case true:
			radiantWins++
		case false:
			direWins++
		}
	}

	radiantRate := (float32(radiantWins) / float32(len(matches))) * 100
	direRate := (float32(direWins) / float32(len(matches))) * 100

	return radiantRate, direRate
}

func formatDuration(totalSeconds int) string {
	minutes := totalSeconds / 60
	seconds := totalSeconds % 60
	return fmt.Sprintf("%d:%02d", minutes, seconds)
}
