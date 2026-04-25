package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/ahmad-alkadri/go-dinorun/internal/app/game"
	"github.com/ahmad-alkadri/go-dinorun/internal/app/scenes"
	"github.com/ahmad-alkadri/go-dinorun/internal/app/sprites"

	"github.com/eiannone/keyboard"
)

const (
	MaxX                int           = 70
	MaxY                int           = 18
	delayCactus         int           = 350
	delayPteranodon     int           = 700
	groundSpeed         int           = 1
	gameSpeed           time.Duration = 15
	delayBetweenEnemies int           = 20
)

func main() {
	if err := keyboard.Open(); err != nil {
		fmt.Println("Failed to open keyboard:", err)
		return
	}
	defer keyboard.Close()

	jumpChan := make(chan struct{}, 1)
	exitChan := make(chan bool, 1)
	var tracker game.SpaceTracker
	go game.HandleInput(jumpChan, exitChan, &tracker)

	for {
		quit := runGameSession(jumpChan, exitChan, &tracker)
		if quit {
			return
		}
		select {
		case <-jumpChan:
		default:
		}
	}
}

// runGameSession plays one round of the game and, on death, blocks on the
// replay prompt. It returns true when the user wants to quit, false to replay.
func runGameSession(jumpChan chan struct{}, exitChan chan bool, tracker *game.SpaceTracker) bool {
	maxX, maxY := MaxX, MaxY
	speed := groundSpeed
	baseY := maxY - 2
	spriteDinoY := baseY
	frameDist := 100

	var (
		dino        sprites.SpriteDino
		cactuses    sprites.SpriteCactuses
		pteranodons sprites.SpritePteranodons
		ground      sprites.SpriteGround
		scores      game.GameScores
		mu          sync.Mutex
	)

	dino.Init(30)
	ground.Init(&maxX)
	scores.Init()

	spawnCactusTicker := time.NewTicker(time.Duration(rand.Intn(1000)+delayCactus) * time.Millisecond)
	spawnPteraTicker := time.NewTicker(time.Duration(rand.Intn(1000)+delayPteranodon) * time.Millisecond)
	defer spawnCactusTicker.Stop()
	defer spawnPteraTicker.Stop()
	defer scores.Stop()

	endGame := func() (finalScene []string, finalScore int) {
		finalScore = scores.Print()
		finalScene = scenes.RenderFinalScene(maxX, maxY, spriteDinoY, speed,
			&dino, &ground, &cactuses, &pteranodons)
		return
	}

	// shouldJump decides whether a SPACE event should trigger a jump right now.
	// The dino can only jump when it is on the ground, and we suppress events
	// that look like OS auto-repeat (the tail of a held key).
	shouldJump := func() bool {
		if spriteDinoY != baseY {
			return false
		}
		return !tracker.IsHeldRepeat()
	}

	for {
		select {
		case <-exitChan:
			finalScene, finalScore := endGame()
			scenes.RenderFinalFrame(finalScene, finalScore)
			return true
		case <-jumpChan:
			if !shouldJump() {
				continue
			}
			T := 12
			maxHeight := 12.0
			displacements := game.GetDisplacements(T, maxHeight)
			for i := 0; i <= 2*T; i++ {
				spriteDinoY -= displacements[i]
				scenes.RenderGame(&maxX, &maxY, &spriteDinoY, &speed,
					&dino, &ground, &cactuses, &pteranodons,
					&scores, exitChan)
				if scenes.AreClashing(&maxY, &spriteDinoY,
					&dino, &cactuses, &pteranodons) {
					finalScene, finalScore := endGame()
					scenes.RenderFinalFrame(finalScene, finalScore)
					return waitForReplay(jumpChan, exitChan, tracker)
				}
				cactuses.Update()
				pteranodons.Update()
				frameDist++
				time.Sleep(gameSpeed * time.Millisecond)
			}
		case <-spawnCactusTicker.C:
			mu.Lock()
			if frameDist > delayBetweenEnemies {
				var newCactus sprites.SpriteCactus
				newCactus.Init(maxX, speed)
				cactuses.Add(newCactus)
				frameDist = 0
			}
			mu.Unlock()
			spawnCactusTicker.Reset(time.Duration(rand.Intn(1000)+delayCactus) * time.Millisecond)
		case <-spawnPteraTicker.C:
			mu.Lock()
			if frameDist > delayBetweenEnemies {
				var newPtera sprites.SpritePteranodon
				newPtera.Init(maxX, speed, 30)
				pteranodons.Add(newPtera)
				frameDist = 0
			}
			mu.Unlock()
			spawnPteraTicker.Reset(time.Duration(rand.Intn(1000)+delayPteranodon) * time.Millisecond)
		default:
			scenes.RenderGame(
				&maxX, &maxY, &spriteDinoY, &speed,
				&dino, &ground, &cactuses, &pteranodons,
				&scores, exitChan)
			if scenes.AreClashing(&maxY, &spriteDinoY,
				&dino, &cactuses, &pteranodons) {
				finalScene, finalScore := endGame()
				scenes.RenderFinalFrame(finalScene, finalScore)
				return waitForReplay(jumpChan, exitChan, tracker)
			}
			cactuses.Update()
			pteranodons.Update()
			frameDist++
			time.Sleep(gameSpeed * time.Millisecond)
		}
	}
}

// waitForReplay blocks until the user makes a fresh SPACE press (replay) or
// asks to quit. A held-space at the moment of death is filtered through the
// same auto-repeat gate the in-game jumps use, so the user must release and
// press again to restart.
func waitForReplay(jumpChan chan struct{}, exitChan chan bool, tracker *game.SpaceTracker) bool {
	for {
		select {
		case <-jumpChan:
			if !tracker.IsHeldRepeat() {
				return false
			}
		case <-exitChan:
			return true
		}
	}
}
