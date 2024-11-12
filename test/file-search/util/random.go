package util

import (
	"database/sql"
	"math/rand"
	"time"
)

func GenerateRandomID() int {
	// Initialize the random number generator with a seed
	rand.Seed(time.Now().UnixNano())

	// Generate a random integer ID (e.g., between 1000 and 9999)
	return rand.Intn(9000) + 1000
}

// NullableString converts a regular string to sql.NullString
func NullableString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{String: "", Valid: false}
	}
	return sql.NullString{String: s, Valid: true}
}
