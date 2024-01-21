package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/eiannone/keyboard"
)

var fruits = []string{"apple", "lemon", "grape", "mango", "peach"}

type Score struct {
	Player int
	AI     int
}

// generateRandomCard: randomly selects a fruit from the slice.
func generateRandomCard(rnd *rand.Rand) string {
	return fruits[rnd.Intn(len(fruits))]
}

// getNewRand: creates a new random number generator.
func getNewRand() *rand.Rand {
	seed := time.Now().UnixNano()
	return rand.New(rand.NewSource(seed))
}

func main() {
	rnd := getNewRand()
	totalRounds := 5
	var score Score

	keyPresses := make(chan rune)
	go handleKeyPresses(keyPresses)

	for round := 1; round <= totalRounds; round++ {
		fmt.Println("------------")
		fmt.Printf("[Round %d] AI:You = %d:%d\n", round, score.AI, score.Player)
		fmt.Println("[s]: start round, [q]: quit game...")
		fmt.Println("[b]: ring the bell if the cards are the same shape!")
		fmt.Println("------------")

		if processRoundInput(rnd, round, &score, keyPresses) {
			break
		}
	}

	fmt.Println("------------")
	fmt.Println("<Game over>")
	fmt.Printf("AI:You = %d:%d\n", score.AI, score.Player)
}

// handleKeyPresses: listens for key presses and sends them to the keyPresses channel.
func handleKeyPresses(keyPresses chan rune) {
	defer keyboard.Close()

	for {
		if char, _, err := keyboard.GetSingleKey(); err == nil {
			keyPresses <- char
		} else {
			return
		}
	}
}

// processRoundInput: handles the input for each round.
func processRoundInput(rnd *rand.Rand, round int, score *Score, keyPresses chan rune) bool {
	for {
		char := <-keyPresses
		switch char {
		case 'q':
			fmt.Println("Quitting game...")
			return true
		case 's':
			playRound(rnd, round, score, keyPresses)
			return false
		default:
			fmt.Println("Invalid input! Try again...")
		}
	}
}

// playRound: manages a single round of the game.
func playRound(rnd *rand.Rand, round int, score *Score, keyPresses chan rune) {
	var playerCard, aiCard string
	var playerRangBell bool

	fmt.Println("Round Start :)")

ROUNDLOOP:
	for {
		playerCard, aiCard = generateRandomCard(rnd), generateRandomCard(rnd)

		fmt.Printf("AI's card: %s\n", aiCard)
		fmt.Printf("Your card: %s\n\n", playerCard)

		delayDuration := calculateWeightedRandomDelay(rnd, round)
		// fmt.Printf("Delaying for %.4f...\n", delayDuration.Seconds())

		select {
		case char := <-keyPresses:
			if char == 'b' {
				playerRangBell = true
				break ROUNDLOOP
			} else if char == 'q' {
				fmt.Println("Quitting game...")
				return
			}
		case <-time.After(delayDuration):
			if playerCard == aiCard {
				fmt.Println("AI rings the bell first! AI wins this round.")
				score.AI++
				return
			}
		}
	}

	if playerRangBell {
		if playerCard == aiCard {
			fmt.Println("You ring the bell! And It's CORRECT! :)")
			fmt.Println("You win this round.")
			score.Player++
		} else {
			fmt.Println("You ring the bell! But It's WRONG. :(")
			fmt.Println("AI wins this round.")
			score.AI++
		}
	}
}

// calculateWeightedRandomDelay: calculates a random delay based on the # of rounds.
func calculateWeightedRandomDelay(rnd *rand.Rand, round int) time.Duration {
	// Start with a fixed maximum delay and decrease it linearly each round
	minDelaySeconds := 0.2
	maxDelaySeconds := 2.0 - float64(round-1)*0.15

	if maxDelaySeconds < 0.5 {
		maxDelaySeconds = 0.5 // Minimum delay of 0.5 seconds
	}

	// Generate a random delay up to maxDelaySeconds
	delayInSeconds := minDelaySeconds + rnd.Float64()*maxDelaySeconds
	return time.Duration(delayInSeconds * float64(time.Second))
}
