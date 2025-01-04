package utils

import (
	"math"

	"github.com/jitendravee/clean_go/internals/models"
)

func CalculateSignalDurations(totalCycle int, signals []models.SingleSignal) []models.SingleSignal {
	minGreen := 10
	minYellow := 3
	maxGreen := totalCycle - minYellow - 2

	totalVehicleCount := 0
	for _, signal := range signals {
		totalVehicleCount += signal.VehicleCount
	}

	// Step 1: Calculate proportional green times
	greenDurations := make([]float64, len(signals))
	for i, signal := range signals {
		proportionalGreen := (float64(signal.VehicleCount) / float64(totalVehicleCount)) * float64(totalCycle)
		greenDurations[i] = math.Max(float64(minGreen), math.Min(proportionalGreen, float64(maxGreen)))
	}

	// Step 2: Adjust green durations if total exceeds available time
	totalGreenTime := 0.0
	for _, green := range greenDurations {
		totalGreenTime += green
	}

	availableGreenTime := float64(totalCycle - len(signals)*minYellow)
	if totalGreenTime > availableGreenTime {
		scalingFactor := availableGreenTime / totalGreenTime
		for i := range greenDurations {
			greenDurations[i] *= scalingFactor
		}
	}

	// Step 3: Assign durations to signals
	for i, signal := range signals {
		signal.GreenDuration = int(greenDurations[i])
		signal.YellowDuration = int(math.Max(float64(minYellow), greenDurations[i]*0.1))
		signal.RedDuration = totalCycle - signal.GreenDuration - signal.YellowDuration
		signals[i] = signal
	}

	return signals
}
