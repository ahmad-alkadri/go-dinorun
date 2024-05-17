package game

func jumpHeight(t, T int, maxY float64) int {
	if t < 0 || t > 2*T {
		return 0 // Outside the interval of interest
	}
	// Compute the height based on the parabolic equation derived and convert to integer.
	return int(-maxY/float64(T*T)*(float64(t-T)*float64(t-T)) + maxY)
}
