package logic

import (
	"math/rand"
	"time"
)

func FraudDetection() (float64, bool) {
	rand.Seed(time.Now().UnixNano())
	fraudIndex := randFloats(0.0, 1.0)
	if fraudIndex > 0.9 {
		return fraudIndex, true
	}

	if fraudIndex <= 0.9 {
		return fraudIndex, false
	}

	return fraudIndex, false
}

func randFloats(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}
