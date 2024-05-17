package sprites

func mapifyStringToScene(strvar string, outerIndex int) map[int]map[int]rune {
	// Create the outer map which contains another map
	result := make(map[int]map[int]rune)
	// Initialize the inner map at key 0
	result[outerIndex] = make(map[int]rune)

	// Loop over each character in the string
	for i, char := range strvar {
		// Add each character to the inner map with its index as the key
		result[outerIndex][i] = char
	}

	return result
}
