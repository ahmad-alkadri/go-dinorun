package game

import (
	"sync"
	"time"
)

// GameScores manages the score of a game, allowing for concurrent updates and
// automatic incrementation over time.
type GameScores struct {
	scores int
	stop   chan struct{}
	mu     sync.Mutex
}

// Init starts a goroutine that increments the score by 1 every 100
// milliseconds. This method should be called to initialize the automatic
// scoring.
func (gs *GameScores) Init() {
	gs.stop = make(chan struct{})
	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				gs.mu.Lock()
				gs.scores++
				gs.mu.Unlock()
			case <-gs.stop:
				return
			}
		}
	}()
}

// Stop halts the automatic score incrementing goroutine. This method should be
// called to clean up the resources when automatic scoring is no longer needed.
func (gs *GameScores) Stop() {
	if gs.stop != nil {
		close(gs.stop)
	}
}

// Add increments the score by a specified value. This method is thread-safe.
func (gs *GameScores) Add(n int) {
	gs.mu.Lock()
	gs.scores += n
	gs.mu.Unlock()
}

// Print returns the current score. This method is thread-safe.
func (gs *GameScores) Print() int {
	gs.mu.Lock()
	defer gs.mu.Unlock()
	return gs.scores
}
