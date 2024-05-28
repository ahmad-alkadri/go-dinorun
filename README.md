# Go Dino Run!

Go Dino Run is a terminal-based implementation of the famous 
[Chrome Dinosaur Game](https://en.wikipedia.org/wiki/Dinosaur_Game),
written in Go.

| ![demo_dark](https://github.com/ahmad-alkadri/go-dinorun/assets/22837764/205a79c1-b214-44c7-a81e-bd8ddeb24214) |
| --- |
| *Game demo, played on terminal, recorded using [asciinema](https://github.com/asciinema)* |

## Installation

You can install Go Dino Run using `go install` directly from GitHub or from a
cloned local repository.

### Install directly from GitHub

```sh
go install github.com/ahmad-alkadri/go-dinorun@latest
```

### Install from a cloned local repository

1. Clone the repository

```sh
git clone https://github.com/ahmad-alkadri/go-dinorun.git
```

2. Navigate to the project directory

```sh
cd go-dinorun
```

3. Install the game

```sh
go install
```

## Uninstallation

To uninstall Go Dino Run, you can remove the installed binary by simply:

```sh
rm -f $(which go-dinorun)
```

## How to Play

Go Dino Run is similar to the Chrome Dinosaur Game. 
The T-Rex will keep running, and you need to avoid the cactuses and pteranodons
by jumping over or ducking below them.

- Open a terminal
- Run `go-dinorun` (make sure the installation's successful)
- Press the `space` button on your keyboard to jump.

Have fun and see how far you can run without hitting the obstacles!

---

## Future Plans

- ~~Adding pteronodons (flying dinosaur)~~
- Adding local database to save some highest scores

Please do not hesitate to raise any issues if you have it, giving me
suggestions, or even contribute.

Enjoy the game!
