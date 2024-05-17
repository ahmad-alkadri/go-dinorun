package sprites

import (
	"math/rand"
	"strings"
)

type SpriteGround struct {
	groundFrame    string
	groundFeatures []string
}

func (sground *SpriteGround) Init(maxX *int, features ...string) {
	if len(features) > 0 {
		sground.groundFeatures = append(sground.groundFeatures, features...)
	} else {
		sground.groundFeatures = []string{"--", "-", "^", "^^", "^-^", "----"}
	}
	sground.groundFrame = sground.generateRandomFeatures(maxX)
}

func (sground *SpriteGround) generateRandomFeatures(maxX *int) (initGroundFrame string) {
	var groundPattern strings.Builder
	desiredLength := *maxX + (len(sground.groundFeatures) * 3) // Slightly over generate to ensure full coverage

	for groundPattern.Len() < desiredLength {
		// Pick a random feature and append it to the pattern
		index := rand.Intn(len(sground.groundFeatures))
		groundPattern.WriteString(sground.groundFeatures[index])
	}

	// Trim the string to exactly maxX characters
	if groundPattern.Len() > *maxX {
		initGroundFrame = groundPattern.String()[:*maxX]
	} else {
		initGroundFrame = groundPattern.String()
	}
	return
}

func (sground *SpriteGround) CurrentGround() string {
	return sground.groundFrame
}

func (sground *SpriteGround) Render(updatedLength int) (renderedground map[int]map[int]rune) {
	// TODO
	// Update the ground by:
	// 1. remove the first three elements
	// 2. generate randomly three elements from the ground pattern
	// 3. append the new three elements at the end
	// Return them as map[int]map[int]rune
	updatedGround := sground.CurrentGround()
	updatedGround = updatedGround[updatedLength:]
	updatedGround += sground.generateRandomFeatures(&updatedLength)
	sground.groundFrame = updatedGround
	return mapifyStringToScene(sground.groundFrame, 0)
}
