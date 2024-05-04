package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/eiannone/keyboard"
)

// TODO
// Change dinoSprite to a type struct
// with a method that return these map-int-rune.
// So instead of it responding to 0 and 1 index, it could create
// some sort of "delay".
// So if the next one is 1 from 0, it'll probably accumulate until 5 "iterations"
// reached before resetting the count and returning different state.
var dinoSprite = []map[int]map[int]rune{
	{
		-3: {2: '\u2584', 3: '\u2584', 4: '\u2584', 5: '\u2584'},              // Head of the dino
		-2: {2: '\u2588', 3: '\u2588', 4: '\u2588', 5: '\u2580'},              // Upper body of the dino
		-1: {0: '\u2580', 1: '\u2588', 2: '\u2588', 3: '\u2588', 4: '\u2588'}, // Lower body of the dino
		0:  {1: 'L', 4: 'L'},                                                  // Feet position 1
	},
	{
		-3: {2: '\u2584', 3: '\u2584', 4: '\u2584', 5: '\u2584'},              // Head of the dino
		-2: {2: '\u2588', 3: '\u2588', 4: '\u2588', 5: '\u2580'},              // Upper body of the dino
		-1: {0: '\u2580', 1: '\u2588', 2: '\u2588', 3: '\u2588', 4: '\u2588'}, // Lower body of the dino
		0:  {2: 'L', 3: 'L'},                                                  // Feet position 2
	},
}

func printSprite(xOffset, yOffset int, sprite map[int]map[int]rune, board []string) {
	for y, row := range sprite {
		for x, char := range row {
			pos := y + yOffset
			if pos >= 0 && pos < len(board) {
				line := []rune(board[pos])
				linePos := x + xOffset
				if linePos >= 0 && linePos < len(line) {
					line[linePos] = char
				}
				board[pos] = string(line)
			}
		}
	}
}

// Turn the whole ground into their own sprite.
var groundOffset int = 0 // Package-level variable, local to this file

func printBoard(MaxX, MaxY, spriteY *int, spriteIndex int) {
	board := make([]string, *MaxY)
	for i := range board {
		board[i] = strings.Repeat(" ", *MaxX)
	}

	// Ground pattern management
	groundPattern := "----____----____"
	repeatCount := (*MaxX / len(groundPattern)) + 2 // Increased repeat count
	fullPattern := strings.Repeat(groundPattern, repeatCount)

	// Calculate the starting position for the ground pattern
	start := groundOffset % len(groundPattern)
	end := start + *MaxX
	if end <= len(fullPattern) {
		board[*MaxY-1] = fullPattern[start:end]
	} else {
		// Wrap around the pattern
		part1 := fullPattern[start:]
		part2 := fullPattern[:end-len(fullPattern)]
		board[*MaxY-1] = part1 + part2
	}

	// Increment ground offset for next frame
	groundOffset++

	// Print the dino sprite
	printSprite(5, *spriteY, dinoSprite[spriteIndex], board)

	// Terminal screen update
	var output strings.Builder
	output.WriteString("\033[H\033[2J\033[3J") // Clear the screen
	for _, line := range board {
		output.WriteString(line + "\n")
	}
	fmt.Print(output.String())
}

func handleInput(jumpChan chan bool, exitChan chan bool) {
	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			keyboard.Close()
			panic(err)
		}
		if key == keyboard.KeySpace || char == ' ' {
			jumpChan <- true
		}
		if key == keyboard.KeyEsc {
			os.Exit(0)
		}
		if key == keyboard.KeyCtrlC {
			fmt.Println("\nExiting...")
			exitChan <- true
			return
		}
	}
}

func main() {
	if err := keyboard.Open(); err != nil {
		fmt.Println("Failed to open keyboard:", err)
		return
	}
	defer func() {
		_ = keyboard.Close()
	}()

	MaxX, MaxY := 60, 20
	baseY := MaxY - 2
	spriteY := baseY
	spriteIndex := 0
	jumpChan := make(chan bool)
	exitChan := make(chan bool)

	go handleInput(jumpChan, exitChan)

	for {
		select {
		case <-exitChan:
			return
		case <-jumpChan:
			if spriteY == baseY {
				// Jump up faster
				for i := 0; i < 6; i++ {
					spriteY -= 2
					printBoard(&MaxX, &MaxY, &spriteY, spriteIndex)
					time.Sleep(50 * time.Millisecond) // Faster jump
				}
				// Fall down slower
				for i := 0; i < 6; i++ {
					spriteY += 2
					printBoard(&MaxX, &MaxY, &spriteY, spriteIndex)
					time.Sleep(50 * time.Millisecond) // Slower fall
				}
			}
		default:
			printBoard(&MaxX, &MaxY, &spriteY, spriteIndex)
			spriteIndex = 1 - spriteIndex
			time.Sleep(50 * time.Millisecond)
		}
	}
}
