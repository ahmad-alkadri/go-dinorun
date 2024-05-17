package game

func GetDisplacements(T int, maxHeight float64) (displacements []int) {
	var arrHeights []int
	for t := 0; t <= 2*T; t++ {
		height := jumpHeight(t, T, maxHeight)
		arrHeights = append(arrHeights, height)
	}
	displacements = calcDisplacements(arrHeights)
	return
}

func calcDisplacements(arrHeights []int) (arrDisps []int) {
	for i, h := range arrHeights {
		if i == 0 {
			arrDisps = append(arrDisps, h)
		} else {
			arrDisps = append(arrDisps, h-arrHeights[i-1])
		}
	}
	return
}
