package main

import (
	"fmt"
	"os"
)

func main() {
	valid := false
	var input int

	for !valid {
		fmt.Println("Get data on the past 100 Pro matches: \n " +
			"1. Get stats for Team \n " +
			"2. Get stats for tournament \n " +
			"3. Get generic stats on Teams \n " +
			"4. Exit")

		fmt.Scan(&input)
		valid = validateInput(input)
	}

	matchData := getMatchData()

	switch input {
	case 1:
		teamStats(matchData)
	case 2:
		tournamentStats(matchData)
	case 3:
		genericStats(matchData)
	case 4:
		os.Exit(0)
	}
}

func validateInput(s int) bool {
	return s > 0 && s <= 4
}
