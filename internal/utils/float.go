package utils

import "math"

func FloatToFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(Round(num * output)) / output
}


func Round(num float64) int64 {
	return int64(num + math.Copysign(0.5, num))
}