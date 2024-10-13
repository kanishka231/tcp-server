package main

import (
	"encoding/json" // Import package for JSON encoding and decoding
	"fmt"           // Import package for formatted I/O
	"net"           // Import package for networking
	"os"            // Import package for operating system functionality
	"strconv"       // Import package for converting strings to integers
	"sync"          // Import package for synchronization primitives like WaitGroup
	"time"          // Import package for time-related functions

	"tcp-server/storage" // Import the storage package that contains Redis-like storage implementation
)

// TPSConfig represents the Transactions Per Second settings from the config file
type TPSConfig struct {
	TPS map[string]int `json:"tps"` // A map to hold TPS values for different time intervals
}

// loadTPSConfig loads the TPS configuration from a JSON file
func loadTPSConfig(filename string) (*TPSConfig, error) {
	data, err := os.ReadFile(filename) // Read the file data
	if err != nil {
		return nil, err // Return error if file reading fails
	}
	var config TPSConfig // Declare a variable to hold the configuration
	err = json.Unmarshal(data, &config) // Unmarshal JSON data into the config variable
	return &config, err // Return the config and any error that occurred
}

func main() {
	// Load TPS config
	tpsConfig, err := loadTPSConfig("./config.json") // Load the TPS configuration from the specified file
	if err != nil {
		fmt.Println("Error loading TPS config:", err) // Print error if loading fails
		return // Exit the program
	}

	// Initialize Redis-like storage
	store := storage.NewRedisLikeStore() // Create a new instance of the Redis-like storage

	// Worker pool to handle sending requests
	var wg sync.WaitGroup // Create a WaitGroup to wait for all goroutines to finish
	for second, tps := range tpsConfig.TPS { // Loop through each second and its TPS
		wg.Add(1) // Increment the WaitGroup counter
		go func(s string, t int) {
			defer wg.Done() // Decrement the counter when the goroutine completes
			sendRequests(s, t, store) // Call the sendRequests function for each second
		}(second, tps)
		time.Sleep(1 * time.Second) // Simulate a delay for each second
	}
	wg.Wait() // Wait for all goroutines to finish
}

// sendRequests sends requests to the TCP server based on the TPS configuration
func sendRequests(_ string, tps int, store *storage.RedisLikeStore) {
	// Number of connections to open
	numConnections := 10 // Define the number of connections to establish
	requestsPerConnection := tps / numConnections // Distribute TPS across connections

	// Limit requests per connection to a maximum of 100
	if requestsPerConnection > 100 { // Check if requests exceed the maximum allowed
		requestsPerConnection = 100 // Set to maximum allowed
	}

	// Create a WaitGroup to wait for all connections to complete
	var wg sync.WaitGroup // Create a WaitGroup for synchronizing goroutines
	for i := 1; i <= numConnections && tps > 0; i++ { // Loop to create connections
		wg.Add(1) // Increment the WaitGroup counter
		go func(p int) {
			defer wg.Done() // Decrement the counter when the goroutine completes

			// Connect to the TCP server
			port := 8000 + p // Calculate the port number to connect to
			conn, err := net.Dial("tcp", "localhost:"+strconv.Itoa(port)) // Establish a TCP connection
			if err != nil {
				fmt.Println("Error connecting to port", port, ":", err) // Log connection error
				return // Exit the goroutine
			}
			defer conn.Close() // Ensure the connection is closed when the function exits

			// Send requests at the TPS rate (limited to requestsPerConnection per connection)
			for j := 0; j < requestsPerConnection && tps > 0; j++ { // Loop to send requests
				// Get data from Redis-like store
				key := fmt.Sprintf("user_%d", j+1) // Create a key for the user
				value, _ := store.Get(key) // Retrieve the value from the storage

				// Send the record to the server
				_, err := conn.Write([]byte(value)) // Send the value to the server
				if err != nil {
					fmt.Println("Error sending data:", err) // Log any sending errors
					return // Exit the goroutine on error
				}

				// Read the server response (transaction ID)
				buf := make([]byte, 1024) // Create a buffer to hold the response
				n, err := conn.Read(buf) // Read the server's response into the buffer
				if err != nil {
					fmt.Println("Error reading response:", err) // Log any reading errors
					return // Exit the goroutine on error
				}
				transactionID := string(buf[:n]) // Convert the response to a string

				// Update Redis-like store with the transaction ID
				store.Set(key, value+" "+transactionID) // Store the updated value with transaction ID

				// Sleep to rate-limit the TPS
				time.Sleep(time.Duration(1000/tps) * time.Millisecond) // Sleep to control TPS
				tps-- // Decrement the TPS for tracking
			}
		}(i) // Pass the current connection number to the goroutine
	}
	wg.Wait() // Wait for all goroutines to finish
}
