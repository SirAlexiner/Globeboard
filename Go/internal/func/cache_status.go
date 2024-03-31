// Package _func provides self defined functions for use in the API.
package _func

import "sync"

var (
	// isCacheCompleted maps language to its cache completion status.
	isCacheCompleted map[string]bool
	// Caching happens in goroutines so initialize mutex.
	mutex sync.Mutex
)

func init() {
	// Initialize the map for cache completion status.
	isCacheCompleted = make(map[string]bool)
}

// GetCacheStatus retrieves the cache completion status for the given language.
func GetCacheStatus(language string) bool {
	return isCacheCompleted[language]
}

// SetCacheStatus sets the cache completion status for the given language.
func SetCacheStatus(language string, status bool) {
	// Mutex lock the cacheCompleted Map to prevent multiple goroutines to write to it concurrently.
	mutex.Lock()
	isCacheCompleted[language] = status
	mutex.Unlock()
	// Unlock the mutex after cache status for language has been updated.
}
