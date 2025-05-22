package main

import (
	"fmt"
	"math/rand"
	"os"
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
		delayCactus         int           = 350
		delayPteranodon     int           = 700
		baseY               int           = MaxY - 2
		spriteDinoY         int           = baseY
		groundSpeed         int           = 1
		gameSpeed           time.Duration = 15
		delayBetweenEnemies int           = 20
		jumpChan            chan bool     = make(chan bool)
		exitChan            chan bool     = make(chan bool)
		dino                sprites.SpriteDino
		cactuses            sprites.SpriteCactuses
		pteranodons         sprites.SpritePteranodons
		ground              sprites.SpriteGround
		spawnCactusTicker   *time.Ticker = time.NewTicker(time.Duration(rand.Intn(1000)+delayCactus) * time.Millisecond)
		spawnPteraTicker    *time.Ticker = time.NewTicker(time.Duration(rand.Intn(1000)+delayPteranodon) * time.Millisecond)
		frameDist           int          = 100
		scores              game.GameScores
		mu                  sync.Mutex
		gameOverScore       int
	)

	// Set up the deferred end sequence
	defer func() {
		scores.Stop()
		gameOverScore = scores.Print()
		game.HandleGameOver(gameOverScore)
		os.Exit(0)
	}()

	dino.Init(30)
	ground.Init(&MaxX)
	scores.Init()

	go game.HandleInput(jumpChan, exitChan)

	for {
		select {
		case <-exitChan:
			gameOverScore = scores.Print()
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
						&scores, exitChan)
					clash := scenes.AreClashing(&MaxY, &spriteDinoY,
						&dino, &cactuses, &pteranodons)
					if clash {
						gameOverScore = scores.Print()
						return
					}
					cactuses.Update()
					pteranodons.Update()
					frameDist += 1
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
			spawnPteraTicker.Reset(time.Duration(rand.Intn(1000)+delayPteranodon) * time.Millisecond)
		default:
			scenes.RenderGame(
				&MaxX, &MaxY, &spriteDinoY, &groundSpeed,
				&dino, &ground, &cactuses, &pteranodons,
				&scores, exitChan)
			clash := scenes.AreClashing(&MaxY, &spriteDinoY,
				&dino, &cactuses, &pteranodons)
			if clash {
				gameOverScore = scores.Print()
				return
			}
			cactuses.Update()
			pteranodons.Update()
			frameDist += 1
			time.Sleep(gameSpeed * time.Millisecond)
		}
	}
}
