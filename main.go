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

func main() {
	if err := keyboard.Open(); err != nil {
		fmt.Println("Failed to open keyboard:", err)
		return
	}
	defer func() {
		_ = keyboard.Close()
	}()
	var (
		MaxX, MaxY          int           = 70, 18
		delayCactus         int           = 1000
		delayPteranodon     int           = 2400
		baseY               int           = MaxY - 2
		spriteDinoY         int           = baseY
		groundSpeed         int           = 1
		gameSpeed           time.Duration = 15
		delayBetweenEnemies int           = 10
		jumpChan            chan bool     = make(chan bool)
		exitChan            chan bool     = make(chan bool)
		gameOverChan        chan bool     = make(chan bool)
	)

	go game.HandleInput(jumpChan, exitChan, gameOverChan)

	// Dino sprite
	var dino sprites.SpriteDino
	dino.Init(30)

	// Cactus sprite
	var cactuses sprites.SpriteCactuses

	// Ptearnodons sprite
	var pteranodons sprites.SpritePteranodons

	var (
		spawnCactusTicker *time.Ticker = time.NewTicker(time.Duration(rand.Intn(1000)+delayCactus) * time.Millisecond)
		spawnPteraTicker  *time.Ticker = time.NewTicker(time.Duration(rand.Intn(1000)+delayPteranodon) * time.Millisecond)
	)

	// Initialize the ground
	var ground sprites.SpriteGround
	ground.Init(&MaxX)

	var (
		// Initialize scores
		scores game.GameScores
		// Distance between enemies
		frameDist int = 100
		mu        sync.Mutex
	)

loop:
	for {
		select {
		case <-exitChan:
			return
		case <-jumpChan:
			if spriteDinoY == baseY {
				T := 12           // Half-duration of the jump, total duration is 2T
				maxHeight := 12.0 // Maximum height of the jump
				displacements := game.GetDisplacements(T, maxHeight)
				for i := 0; i <= 2*T; i++ {
					spriteDinoY -= displacements[i]
					scenes.RenderGame(&MaxX, &MaxY, &spriteDinoY, &groundSpeed,
						&dino, &ground, &cactuses, &pteranodons,
						&scores, gameOverChan)
					clash := scenes.AreClashing(&MaxY, &spriteDinoY,
						&dino, &cactuses, &pteranodons)
					if clash {
						break loop
					}
					// Check score based on passed cactuses
					diff, deltaX, _ := cactuses.Update()
					if diff > 0 {
						if deltaX == 4 {
							scores.Add(1)
						}
						if deltaX == 10 {
							scores.Add(2)
						}
					}
					// Check score based on passed pteranodons
					diffP, _, _ := pteranodons.Update()
					if diffP > 0 {
						scores.Add(1)
					}
					// Add framediff
					frameDist += 1
					// Time sleep before next frame
					time.Sleep(gameSpeed * time.Millisecond)
				}
			}
		case <-spawnCactusTicker.C:
			mu.Lock()
			if frameDist > delayBetweenEnemies {
				var newCactus sprites.SpriteCactus
				newCactus.Init(MaxX, groundSpeed)
				cactuses.Add(newCactus)
				frameDist = 0
			}
			mu.Unlock()
			// Reset ticker
			spawnCactusTicker.Reset(time.Duration(rand.Intn(1000)+delayCactus) * time.Millisecond)
		case <-spawnPteraTicker.C:
			mu.Lock()
			if frameDist > delayBetweenEnemies {
				var newPtera sprites.SpritePteranodon
				newPtera.Init(MaxX, groundSpeed, 30)
				pteranodons.Add(newPtera)
				frameDist = 0
			}
			mu.Unlock()
			// Reset ticker
			spawnPteraTicker.Reset(time.Duration(rand.Intn(1000)+delayPteranodon) * time.Millisecond)
		default:
			scenes.RenderGame(
				&MaxX, &MaxY, &spriteDinoY, &groundSpeed,
				&dino, &ground, &cactuses, &pteranodons,
				&scores, gameOverChan)
			clash := scenes.AreClashing(&MaxY, &spriteDinoY,
				&dino, &cactuses, &pteranodons)
			if clash {
				break loop
			}
			// Check score based on passed cactuses
			diff, deltaX, _ := cactuses.Update()
			if diff > 0 {
				if deltaX == 4 {
					scores.Add(1)
				}
				if deltaX == 10 {
					scores.Add(2)
				}
			}
			// Check score based on passed pteranodons
			diffP, _, _ := pteranodons.Update()
			if diffP > 0 {
				scores.Add(1)
			}
			// Add framediff
			frameDist += 1
			// Time sleep before next frame
			time.Sleep(gameSpeed * time.Millisecond)
		}
	}

	go func() {
		gameOverChan <- true
	}()
}
