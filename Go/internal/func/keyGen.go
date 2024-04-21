// Package _func provides developer-made utility functions for use within the application.
package _func

import (
	"math/rand"
	"time"
)

// init initializes the package by seeding the random number generator with the current time.
func init() {
	rand.New(rand.NewSource(time.Now().UnixNano())) // Seed rand with time, making it nondeterministic.
}

// Runes is the list of characters to use for generating IDs.
var Runes = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")

// GenerateAPIKey generates an API key of length 'n'.
func GenerateAPIKey(n int) string {
	// Make two slices that are 'n' long, a and b.
	a := make([]rune, n)
	b := make([]rune, n)
	// loop through the slices and insert a random character at each index for each slice.
	for i := range a {
		a[i] = Runes[rand.Intn(len(Runes))]
		b[i] = Runes[rand.Intn(len(Runes))]
	}
	// concatenate the slices to strings, following standard format.
	return "sk-" + string(a) + "tTRjPv" + string(b)
}

// GenerateUID returns a Unique Identifier of length 'n'.
func GenerateUID(n int) string {
	// Make a slice that is 'n' long.
	b := make([]rune, n)
	// loop through the slice and insert a random character at each index.
	for i := range b {
		b[i] = Runes[rand.Intn(len(Runes))]
	}
	// concatenate the slice to string.
	return string(b)
}
