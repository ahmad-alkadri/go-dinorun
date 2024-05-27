package sprites

import (
	"math"
	"math/rand"
)

type SpritePteranodons struct {
	Group []SpritePteranodon
}

func (spteras *SpritePteranodons) Update() (diff, deltaX, deltaY int) {
	var activePteras []SpritePteranodon
	var disappearingPtera SpritePteranodon
	for _, ptera := range spteras.Group {
		ptera.UpdatePosition()
		if ptera.Xoffset > 0 { // Check if still on screen
			activePteras = append(activePteras, ptera)
		} else {
			disappearingPtera = ptera
		}
	}
	diff = int(math.Abs(float64(len(spteras.Group) - len(activePteras))))
	spteras.Group = activePteras
	deltaX, deltaY = disappearingPtera.SpanCells()
	return
}

func (spteras *SpritePteranodons) Add(sptera SpritePteranodon) {
	spteras.Group = append(spteras.Group, sptera)
}

type SpritePteranodon struct {
	frameCount int
	frameLimit int
	halfLim    float32
	Xoffset    int
	Yoffset    int
	Graphic    map[int]map[int]rune
	deltaXRate int
}

func (sptera *SpritePteranodon) Init(MaxX int, deltaXRate int, frameLimit int) {
	// Make sure the framecount starts from zero
	sptera.resetFrameCount()
	// Make sure that the minimum frame limit is 2
	if frameLimit < 2 {
		sptera.frameLimit = 2
	} else {
		sptera.frameLimit = frameLimit
	}
	sptera.frameLimit = 30
	sptera.halfLim = float32(frameLimit/2) - 1
	sptera.Xoffset = MaxX
	sptera.Yoffset = sptera.randomYoffset()
	sptera.deltaXRate = deltaXRate
}

func (spetra *SpritePteranodon) SpanCells() (deltaX, deltaY int) {
	// Initialize min and max values
	var minX, minY int = math.MaxInt, math.MaxInt
	var maxX, maxY int = math.MinInt, math.MinInt

	// Iterate over the slice of maps and find min and max values
	for y, innerMap := range spetra.Graphic {
		for x := range innerMap {
			if x < minX {
				minX = x
			}
			if x > maxX {
				maxX = x
			}
			if y < minY {
				minY = y
			}
			if y > maxY {
				maxY = y
			}
		}
	}
	deltaX = maxX - minX
	deltaY = maxY - minY
	return
}

func (sptera *SpritePteranodon) resetFrameCount() {
	sptera.frameCount = 0
}

func (sptera *SpritePteranodon) randomYoffset() int {
	return rand.Intn(6) - 10
}

func (sptera *SpritePteranodon) increment() {
	if sptera.frameCount == sptera.frameLimit {
		sptera.resetFrameCount()
	}
	sptera.frameCount += 1
	// Decide which graph to be put now
	// Check its value
	if sptera.firstHalfOrNot() {
		// Return the first half frame
		sptera.Graphic = sptera.frame(0)
	} else {
		// Return the second half frame
		sptera.Graphic = sptera.frame(1)
	}
}

func (sptera *SpritePteranodon) firstHalfOrNot() bool {
	// Get the limit and halven it
	halfLim := float32(sptera.frameLimit)/2 - 1
	// Check the current frameCount
	return float32(sptera.frameCount) <= halfLim
}

func (sptera *SpritePteranodon) Render() map[int]map[int]rune {
	// Return the graphic
	return sptera.Graphic
}

func (sptera *SpritePteranodon) UpdatePosition() {
	// Increase framecount
	sptera.increment()
	sptera.Xoffset -= sptera.deltaXRate
}

func (sptera *SpritePteranodon) frame(i int) map[int]map[int]rune {
	dinoSprite := [2]map[int]map[int]rune{
		{
			sptera.Yoffset - 1: {0: '\u2597', 1: '\u2584', 3: '\u2588'},
			sptera.Yoffset:     {0: '\u2580', 1: '\u2588', 2: '\u2588', 3: '\u2588', 4: '\u2588', 5: '\u2584'},
		},
		{
			sptera.Yoffset - 1: {0: '\u2597', 1: '\u2584'},
			sptera.Yoffset:     {0: '\u2580', 1: '\u2588', 2: '\u2588', 3: '\u2588', 4: '\u2588', 5: '\u2584'},
			sptera.Yoffset + 1: {3: '\u2580'},
		},
	}
	return dinoSprite[i]
}
