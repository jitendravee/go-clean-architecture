package utils

import "math"

func CalculateSignalDurations(vehicleCount, totalCycle int) (int, int, int) {

	minGreen := 10
	maxGreen := 60
	yellow := 5

	green := float64(minGreen) + float64(vehicleCount)
	if green > float64(maxGreen) {
		green = float64(maxGreen)
	}

	red := float64(totalCycle) - (green + float64(yellow))
	if red < 0 {
		red = 0
	}

	return int(math.Round(green)), yellow, int(math.Round(red))
}
