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
	gameOverChan chan bool,
	sceneBuffer []string,
) {
	// Reuse the scene buffer to avoid reallocations
	for i := range sceneBuffer {
		sceneBuffer[i] = strings.Repeat(" ", *MaxX)
	}

	// Print the dino sprite
	printSprite(5, *spriteDinoY, dino.Render(), sceneBuffer)

	// Print the ground as sprite
	printSprite(0, *MaxY-1, ground.Render(*groundSpeed), sceneBuffer)

	// Print the cactuses
	for _, cactus := range cactuses.Group {
		cactusXoffset := cactus.Xoffset
		printSprite(cactusXoffset, *MaxY-1, cactus.Render(), sceneBuffer)
	}

	// Print the pteranodons
	for _, ptera := range pteranodons.Group {
		pteraXoffset := ptera.Xoffset
		printSprite(pteraXoffset, *MaxY-1, ptera.Render(), sceneBuffer)
	}

	// Terminal screen update
	// Using a single write to avoid multiple I/O operations
	var output strings.Builder
	output.Grow((*MaxY + 2) * (*MaxX))         // Preallocate enough space
	output.WriteString("\033[H\033[2J\033[3J") // Clear the screen
	output.WriteString(fmt.Sprintf("Score: %d\n", scores.Print()))
	for _, line := range sceneBuffer {
		output.WriteString(line)
		output.WriteByte('\n')
	}
	fmt.Print(output.String())
}

// Helper function to initialize the scene buffer
func InitializeSceneBuffer(MaxY, MaxX int) []string {
	sceneBuffer := make([]string, MaxY)
	for i := range sceneBuffer {
		sceneBuffer[i] = strings.Repeat(" ", MaxX)
	}
	return sceneBuffer
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
