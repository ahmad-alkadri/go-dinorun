package sprites

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GroundLength(t *testing.T) {
	// Initialize the ground
	var sground SpriteGround
	var maxX int = 10
	sground.Init(&maxX)
	curground := sground.CurrentGround()

	fmt.Printf("Current ground: %s", curground)
	// fmt.Println(curground)

	// Assert the length
	assert.Equal(t, maxX, len(curground),
		"Ground length not fulfilled. Expected %d, found %d",
		maxX, len(curground))
}
