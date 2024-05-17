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
	scores *game.GameScores,
	gameOverChan chan bool,
) {
	scene := make([]string, *MaxY)
	for i := range scene {
		scene[i] = strings.Repeat(" ", *MaxX)
	}

	// Print the dino sprite
	dinoFig := dino.Render()
	printSprite(5, *spriteDinoY, dinoFig, scene)

	// Print the ground as sprite
	printSprite(0, *MaxY-1, ground.Render(*groundSpeed), scene)

	// Print the cactuses
	for _, cactus := range cactuses.Group {
		cactusXoffset := cactus.Xoffset
		printSprite(cactusXoffset, *MaxY-1, cactus.Render(), scene)
	}

	// Terminal screen update
	var output strings.Builder
	output.WriteString("\033[H\033[2J\033[3J") // Clear the screen
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
) bool {
	// Check if there's any clash
	dinoCells := extractSpriteCells(5, *spriteDinoY, dino.Render())
	var cactusCells [][2]int
	for _, cactus := range cactuses.Group {
		cactusCells = extractSpriteCells(cactus.Xoffset, *MaxY-1, cactus.Render())
		if shareChild(dinoCells, cactusCells) {
			return true
		}
	}
	return false
}
