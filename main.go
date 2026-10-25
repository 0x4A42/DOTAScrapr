package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	cache := NewMatchCache(5 * time.Minute)

	for {
		fmt.Println("For the past 100 Pro matches: \n " +
			"1. Get stats for Team \n " +
			"2. Get stats for tournament \n " +
			"3. Get generic stats on Teams \n " +
			"4. Exit")

		var input int
		_, err := fmt.Scan(&input)
		if err != nil || !validateInput(input) {
			fmt.Println("Invalid input. Please choose an option between 1 and 4.")
			var discard string
			_, _ = fmt.Scanln(&discard)
			continue
		}

		if input == 4 {
			fmt.Println("Exiting...")
			os.Exit(0)
		}

		matchData, err := cache.GetData()
		if err != nil {
			fmt.Printf("Failed to get match data: %v\n", err)
			return
		}

		switch input {
		case 1:
			teamStats(matchData)
		case 2:
			tournamentStats(matchData)
		case 3:
			genericStats(matchData)
		}

		fmt.Println()
	}
}

func validateInput(s int) bool {
	return s > 0 && s <= 4
}
