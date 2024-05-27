package sprites

import (
	"math/rand"
)

var PossibleForms = []map[int]map[int]rune{
	{
		-4: {2: '\u2588', 8: '\u2588'},
		-3: {2: '\u2588', 4: '\u2588', 8: '\u2588', 10: '\u2588'},
		-2: {0: '\u2588', 2: '\u2588', 3: '\u2580', 4: '\u2580', 6: '\u2588', 8: '\u2588', 9: '\u2580', 10: '\u2580'},
		-1: {0: '\u2580', 1: '\u2580', 2: '\u2588', 6: '\u2580', 7: '\u2580', 8: '\u2588'},
		0:  {2: '\u2580', 8: '\u2580'},
	},
	{
		-4: {2: '\u2588'},
		-3: {2: '\u2588', 4: '\u2588'},
		-2: {0: '\u2588', 2: '\u2588', 3: '\u2580', 4: '\u2580'},
		-1: {0: '\u2580', 1: '\u2580', 2: '\u2588'},
		0:  {2: '\u2580'},
	},
	{
		-3: {2: '\u2588', 8: '\u2588'},
		-2: {0: '\u2588', 2: '\u2588', 4: '\u2588', 6: '\u2588', 8: '\u2588', 10: '\u2588'},
		-1: {0: '\u2580', 1: '\u2580', 2: '\u2588', 3: '\u2580', 4: '\u2580', 6: '\u2580', 7: '\u2580', 8: '\u2588', 9: '\u2580', 10: '\u2580'},
		0:  {2: '\u2580', 8: '\u2580'},
	},
	{
		-3: {2: '\u2588'},
		-2: {0: '\u2588', 2: '\u2588', 4: '\u2588'},
		-1: {0: '\u2580', 1: '\u2580', 2: '\u2588', 3: '\u2580', 4: '\u2580'},
		0:  {2: '\u2580'},
	},
}

type SpriteCactuses struct {
	Group []SpriteCactus
}

func (scacts *SpriteCactuses) Update() {
	var activeCactuses []SpriteCactus
	for _, cactus := range scacts.Group {
		cactus.UpdatePosition()
		if cactus.Xoffset > 0 { // Check if still on screen
			activeCactuses = append(activeCactuses, cactus)
		}
	}
	scacts.Group = activeCactuses
}

func (scacts *SpriteCactuses) Add(scact SpriteCactus) {
	scacts.Group = append(scacts.Group, scact)
}

type SpriteCactus struct {
	Xoffset    int
	Graphic    map[int]map[int]rune
	deltaXRate int
}

func (scact *SpriteCactus) Init(MaxX int, deltaXRate ...int) {
	scact.Graphic = scact.forms()
	scact.Xoffset = MaxX
	// Delta X rate by default
	scact.deltaXRate = 3
	if len(deltaXRate) > 0 {
		scact.deltaXRate = deltaXRate[0]
	}
}

func (scact *SpriteCactus) Render() map[int]map[int]rune {
	return scact.Graphic
}

func (scact *SpriteCactus) UpdatePosition() {
	scact.Xoffset -= scact.deltaXRate
}

func (scact *SpriteCactus) forms() map[int]map[int]rune {
	// Choose one random form to be rendered for this cactus
	idx := scact.choseRandomForm(len(PossibleForms))
	return PossibleForms[idx]
}

func (scact *SpriteCactus) choseRandomForm(qtyOfForms int) int {
	// Return a random integer in the range [0, n].
	if qtyOfForms == 1 {
		return 0
	}
	return rand.Intn(qtyOfForms)
}
