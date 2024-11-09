package util

import (
	"math/rand"
	"time"
)

func GenerateRandomID() int {
	// Initialize the random number generator with a seed
	rand.Seed(time.Now().UnixNano())

	// Generate a random integer ID (e.g., between 1000 and 9999)
	return rand.Intn(9000) + 1000
}
