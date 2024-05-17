package game

import (
	"fmt"
	"os"

	"github.com/eiannone/keyboard"
)

func HandleInput(jumpChan chan bool, exitChan chan bool, gameOverChan chan bool) {
	
	for {
		char, key, _ := keyboard.GetKey()

		defer func() {
			_ = keyboard.Close()
		}()

		select {
        case <-gameOverChan:
			fmt.Println("\nGame Over...")
			exitChan <- true
			return
        default:
			if key == keyboard.KeySpace || char == ' ' {
				jumpChan <- true
			}
			if key == keyboard.KeyEsc {
				os.Exit(0)
			}
			if key == keyboard.KeyCtrlC {
				fmt.Println("\nExiting...")
				exitChan <- true
				return
			}
        }
	}
}
