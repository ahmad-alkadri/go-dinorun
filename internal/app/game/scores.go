package game

type GameScores struct {
	scores int
}

func (gs *GameScores) Add(n int) {
	gs.scores += n
}

func (gs *GameScores) Print() int {
	return gs.scores
}
