package game

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/eiannone/keyboard"
)

func HandleInput(jumpChan chan bool, exitChan chan bool) {
	for {
		char, key, _ := keyboard.GetKey()

		if key == keyboard.KeySpace || char == ' ' {
			jumpChan <- true
		}
		if key == keyboard.KeyEsc || key == keyboard.KeyCtrlC {
			fmt.Print("\nExiting...\n")
			exitChan <- true
			return
		}

		select {
		case <-exitChan:
			return
		default:
		}
	}
}

// HandleGameOver manages the game over sequence including high score handling
func HandleGameOver(score int) {
	fmt.Printf("\nGame Over!\nFinal Score: %d\n", score)

	// Check if this is a high score
	if IsHighScore(score) {
		// Get player name and save high score
		name := GetPlayerName()
		if err := AddHighScore(name, score); err != nil {
			fmt.Printf("\nError saving high score: %v\n", err)
		} else {
			// Display the high scores
			if scores, err := LoadHighScores(); err == nil {
				fmt.Println("\nHigh Scores:")
				fmt.Println("--------------------")
				for i, score := range scores {
					fmt.Printf("%d. %-20s %d\n", i+1, score.Name, score.Score)
				}
				fmt.Println("--------------------")
				return
			}
		}
	}
}

// GetPlayerName prompts for and returns the player's name
func GetPlayerName() string {
	fmt.Print("\nCongratulations! You got a high score!\nEnter your name: ")
	reader := bufio.NewReader(os.Stdin)
	name, _ := reader.ReadString('\n')
	return strings.TrimSpace(name)
}
