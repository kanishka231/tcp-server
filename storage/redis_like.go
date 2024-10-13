package storage // Declare the package name

import (
	"fmt"       // Import package for formatted I/O
	"sync"      // Import package for synchronization primitives like RWMutex
)

// RedisLikeStore represents a simple in-memory key-value store resembling Redis
type RedisLikeStore struct {
	store map[string]string // A map to store key-value pairs
	mu    sync.RWMutex      // A mutex for handling concurrent access
}

// NewRedisLikeStore initializes the Redis-like store with 500,000 records
func NewRedisLikeStore() *RedisLikeStore {
	r := &RedisLikeStore{ // Create a new instance of RedisLikeStore
		store: make(map[string]string), // Initialize the store as an empty map
	}
	// Simulating 500,000 records
	for i := 1; i <= 500000; i++ { // Loop to create 500,000 records
		key := fmt.Sprintf("user_%d", i) // Generate a unique key for each user
		value := fmt.Sprintf(`{"name": "user_%d", "age": %d}`, i, 20+(i%30)) // Create a JSON-like string as the value
		r.store[key] = value // Store the key-value pair in the store
	}
	return r // Return the initialized store
}

// Get retrieves a record by key from the store
func (r *RedisLikeStore) Get(key string) (string, bool) {
	r.mu.RLock() // Acquire a read lock for concurrent access
	defer r.mu.RUnlock() // Ensure the lock is released when the function exits
	value, exists := r.store[key] // Retrieve the value and check for existence
	return value, exists // Return the value and a boolean indicating existence
}

// Set updates a record by key in the store (used for adding transaction_id from the server)
func (r *RedisLikeStore) Set(key string, value string) {
	r.mu.Lock() // Acquire a write lock for concurrent access
	defer r.mu.Unlock() // Ensure the lock is released when the function exits
	r.store[key] = value // Update the store with the new value for the given key
}
