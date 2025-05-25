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
	// Terminal screen update with enhanced buffer control
	var output strings.Builder
	output.WriteString("\u001B[?1049h") // Enable alternate screen buffer
	output.WriteString("\u001B[?25l")   // Hide cursor
	output.WriteString("\u001B[2J")     // Clear screen
	output.WriteString("\u001B[H")      // Move cursor to home position
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

// RenderFinalFrame renders the last frame of the game with the game over message
// while maintaining the alternate buffer to prevent scrollback
func RenderFinalFrame(scene []string, score int) {
	var output strings.Builder
	// Stay in alternate buffer, just clear and redraw
	output.WriteString("\u001B[2J")   // Clear screen
	output.WriteString("\u001B[H")    // Move cursor to home position
	output.WriteString("\u001B[?25h") // Show cursor

	// Print the final scene
	output.WriteString(fmt.Sprintf("Score: %d\n", score))
	for _, line := range scene {
		output.WriteString(line + "\n")
	}

	// Print game over message below the frame
	output.WriteString("\nGame Over!\n")
	output.WriteString("Press any key to exit...")
	fmt.Print(output.String())
}

// RenderFinalScene creates the final scene with all game elements in their last positions
func RenderFinalScene(maxX, maxY, spriteDinoY, groundSpeed int,
	dino *sprites.SpriteDino,
	ground *sprites.SpriteGround,
	cactuses *sprites.SpriteCactuses,
	pteranodons *sprites.SpritePteranodons) []string {

	finalScene := make([]string, maxY)
	for i := range finalScene {
		finalScene[i] = strings.Repeat(" ", maxX)
	}

	// Print the dino sprite
	printSprite(5, spriteDinoY, dino.Render(), finalScene)

	// Print the ground as sprite
	printSprite(0, maxY-1, ground.Render(groundSpeed), finalScene)

	// Print the cactuses
	for _, cactus := range cactuses.Group {
		printSprite(cactus.Xoffset, maxY-1, cactus.Render(), finalScene)
	}

	// Print the pteranodons
	for _, ptera := range pteranodons.Group {
		printSprite(ptera.Xoffset, maxY-1, ptera.Render(), finalScene)
	}

	return finalScene
}
