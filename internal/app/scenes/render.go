package scenes

import (
	"fmt"
	"strings"

	"github.com/ahmad-alkadri/go-dinorun/internal/app/game"
	"github.com/ahmad-alkadri/go-dinorun/internal/app/sprites"
)

func RenderGame(
	MaxX, MaxY, spriteDinoY, groundSpeed *int,
	dino *sprites.SpriteDino,
	ground *sprites.SpriteGround,
	cactuses *sprites.SpriteCactuses,
	pteranodons *sprites.SpritePteranodons,
	scores *game.GameScores,
	exitChan chan bool,
) {
	scene := make([]string, *MaxY)
	for i := range scene {
		scene[i] = strings.Repeat(" ", *MaxX)
	}

	// Print the dino sprite
	printSprite(5, *spriteDinoY, dino.Render(), scene)

	// Print the ground as sprite
	printSprite(0, *MaxY-1, ground.Render(*groundSpeed), scene)

	// Print the cactuses
	for _, cactus := range cactuses.Group {
		cactusXoffset := cactus.Xoffset
		printSprite(cactusXoffset, *MaxY-1, cactus.Render(), scene)
	}

	// Print the pteranodons
	for _, ptera := range pteranodons.Group {
		pteraXoffset := ptera.Xoffset
		printSprite(pteraXoffset, *MaxY-1, ptera.Render(), scene)
	}

	// Terminal screen update
	var output strings.Builder
	// Using PowerShell-friendly clear screen approach
	output.WriteString("\u001B[2J\u001B[H") // Clear screen and move cursor to home position
	output.WriteString(fmt.Sprintf("Score: %d\n", scores.Print()))
	for _, line := range scene {
		output.WriteString(line + "\n")
	}
	fmt.Print(output.String())
}

func AreClashing(
	MaxY, spriteDinoY *int,
	dino *sprites.SpriteDino,
	cactuses *sprites.SpriteCactuses,
	pteranodons *sprites.SpritePteranodons,
) bool {
	// Check if there's any clash
	dinoCells := extractSpriteCells(5, *spriteDinoY, dino.Render())
	var cactusCells [][2]int
	var pteraCells [][2]int
	// Clash with the cactus
	for _, cactus := range cactuses.Group {
		cactusCells = extractSpriteCells(cactus.Xoffset, *MaxY-1, cactus.Render())
		if shareChild(dinoCells, cactusCells) {
			return true
		}
	}
	// Clash with the pteranodons
	for _, ptera := range pteranodons.Group {
		pteraCells = extractSpriteCells(ptera.Xoffset, *MaxY-1, ptera.Render())
		if shareChild(dinoCells, pteraCells) {
			return true
		}
	}
	return false
}
