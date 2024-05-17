package sprites

// The SpriteDino type in Go represents a sprite animation of a dinosaur character with different
// frames for rendering.
//
// @property {int} frameCount - The `frameCount` property in the `SpriteDino` struct represents the
// current frame count of the sprite animation. It keeps track of how many frames have been rendered so
// far.
//
// @property {int} frameLimit - The `frameLimit` property in the `SpriteDino` struct represents the
// maximum number of frames allowed before resetting the frame count back to 0. This property is used
// to control the animation of the sprite.
type SpriteDino struct {
	frameCount int
	frameLimit int
}

func (sdino *SpriteDino) Init(frameLimit ...int) {
	sdino.resetFrameCount()
	if len(frameLimit) > 0 {
		sdino.frameLimit = frameLimit[0]
	} else {
		sdino.frameLimit = 2
	}
}

func (sdino *SpriteDino) add() {
	sdino.frameCount++
}

func (sdino *SpriteDino) resetFrameCount() {
	sdino.frameCount = 0
}

func (sdino *SpriteDino) increment() {
	if sdino.frameCount == sdino.frameLimit {
		sdino.resetFrameCount()
	}
	sdino.add()
}

func (sdino *SpriteDino) firstHalfOrNot() bool {
	// Get the limit and halven it
	halfLim := float32(sdino.frameLimit)/2-1
	// Check the current frameCount
	return float32(sdino.frameCount) <= halfLim
}

func (sdino *SpriteDino) Render() map[int]map[int]rune {
	// Increase framecount
	sdino.increment()
	// Check its value
	if sdino.firstHalfOrNot() {
		// Return the first half frame
		return sdino.frame(0)
	} else {
		// Return the second half frame
		return sdino.frame(1)
	}
}

func (sdino *SpriteDino) frame(i int) map[int]map[int]rune {
	dinoSprite := [2]map[int]map[int]rune{
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
	return dinoSprite[i]
}