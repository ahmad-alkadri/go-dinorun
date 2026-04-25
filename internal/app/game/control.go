package game

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/eiannone/keyboard"
)

// AutoRepeatGap is the maximum interval between two consecutive SPACE events
// that we treat as part of the same physical key-hold (OS auto-repeat fires
// well below this, while a deliberate release-and-press cycle takes longer).
// The game loop uses this to suppress the auto-repeat tail of a held key.
const AutoRepeatGap = 100 * time.Millisecond

// SpaceTracker records the timestamps of the two most recent SPACE events.
// The input goroutine writes; the game loop reads. Atomics keep both sides
// lock-free; the tracker is independent from the trigger channel so that
// dropped channel sends never lose timing information.
type SpaceTracker struct {
	last atomic.Int64 // unix nanos of most recent event, 0 if never pressed
	prev atomic.Int64 // unix nanos of the event before that, 0 if only one
}

func (s *SpaceTracker) record(now time.Time) {
	prev := s.last.Load()
	s.last.Store(now.UnixNano())
	s.prev.Store(prev)
}

// RecentGap returns the interval between the two most recent SPACE events,
// or 0 if fewer than two events have been recorded. A small non-zero gap
// indicates the key is being held (OS auto-repeat).
func (s *SpaceTracker) RecentGap() time.Duration {
	last := s.last.Load()
	prev := s.prev.Load()
	if last == 0 || prev == 0 {
		return 0
	}
	return time.Duration(last - prev)
}

// IsHeldRepeat reports whether the latest event looks like a continuation
// of a held key (recent gap below the auto-repeat threshold).
func (s *SpaceTracker) IsHeldRepeat() bool {
	gap := s.RecentGap()
	return gap > 0 && gap < AutoRepeatGap
}

func HandleInput(jumpChan chan struct{}, exitChan chan bool, tracker *SpaceTracker) {
	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			select {
			case exitChan <- true:
			default:
			}
			return
		}

		if key == keyboard.KeySpace || char == ' ' {
			tracker.record(time.Now())
			// Non-blocking send: the channel is just a wakeup signal. If the
			// game loop hasn't drained yet, the tracker still captured the
			// timing, so dropping the send is safe.
			select {
			case jumpChan <- struct{}{}:
			default:
			}
		}
		if key == keyboard.KeyEsc || key == keyboard.KeyCtrlC {
			fmt.Print("\nExiting...\n")
			exitChan <- true
			return
		}

		select {
		case <-exitChan:
			return
		default:
		}
	}
}
