package scenes

func printSprite(xOffset, yOffset int, sprite map[int]map[int]rune, scene []string) {
	for y, row := range sprite {
		for x, char := range row {
			pos := y + yOffset
			if pos >= 0 && pos < len(scene) {
				line := []rune(scene[pos])
				linePos := x + xOffset
				if linePos >= 0 && linePos < len(line) {
					line[linePos] = char
				}
				scene[pos] = string(line)
			}
		}
	}
}

func extractSpriteCells(
	xOffset, yOffset int,
	sprite map[int]map[int]rune,
) (cells [][2]int) {
	for y, row := range sprite {
		for x := range row {
			pos := y + yOffset
			if pos >= 0 {
				linePos := x + xOffset
				cell := [2]int{linePos, pos}
				cells = append(cells, cell)
			}
		}
	}
	return
}

func shareChild(cells1, cells2 [][2]int) bool {
	for _, pair1 := range cells1 {
		for _, pair2 := range cells2 {
			if pair1[0] == pair2[0] && pair1[1] == pair2[1] {
				return true
			}
		}
	}
	return false
}
