package main

import (
	"fmt"
	"math/rand"
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

	MaxX, MaxY := 70, 18
	baseY := MaxY - 2
	spriteDinoY := baseY
	groundSpeed := 1

	jumpChan := make(chan bool)
	exitChan := make(chan bool)
	gameOverChan := make(chan bool)

	go game.HandleInput(jumpChan, exitChan, gameOverChan)

	// INITIALIZING THE SPRITES
	// Dino sprite
	var dino sprites.SpriteDino
	dino.Init(30)

	// Cactus sprite
	var cactuses sprites.SpriteCactuses

	spawnCactusTicker := time.NewTicker(time.Duration(rand.Intn(1000)+350) * time.Millisecond)

	// Initialize the ground
	var ground sprites.SpriteGround
	ground.Init(&MaxX)

	// Initialize scores
	var scores game.GameScores

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
						&dino, &ground, &cactuses,
						&scores, gameOverChan)
					clash := scenes.AreClashing(&MaxY, &spriteDinoY,
						&dino, &cactuses)
					if clash {
						break loop
					}
					diff, deltaX, _ := cactuses.Update()
					if diff > 0 {
						if deltaX == 4 {
							scores.Add(1)
						}
						if deltaX == 10 {
							scores.Add(2)
						}
					}
					time.Sleep(15 * time.Millisecond)
				}
			}
		case <-spawnCactusTicker.C:
			var newCactus sprites.SpriteCactus
			newCactus.Init(MaxX, groundSpeed)
			cactuses.Add(newCactus)
			// Reset ticker
			spawnCactusTicker.Reset(time.Duration(rand.Intn(1000)+350) * time.Millisecond)
		default:
			scenes.RenderGame(
				&MaxX, &MaxY, &spriteDinoY, &groundSpeed,
				&dino, &ground, &cactuses,
				&scores, gameOverChan)
			clash := scenes.AreClashing(&MaxY, &spriteDinoY,
				&dino, &cactuses)
			if clash {
				break loop
			}
			diff, deltaX, _ := cactuses.Update()
			if diff > 0 {
				if deltaX == 5 {
					scores.Add(1)
				}
				if deltaX == 11 {
					scores.Add(2)
				}
			}
			time.Sleep(15 * time.Millisecond)
		}
	}

	go func() {
		gameOverChan <- true
	}()
}
