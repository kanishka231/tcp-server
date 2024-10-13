package main

import (
	"fmt"      // Import the fmt package for formatted I/O
	"net"      // Import the net package for network operations
	"strconv"  // Import the strconv package for converting strings to integers
	"time"     // Import the time package for working with time-related functions
)

// handleConnection handles incoming TCP connections.
func handleConnection(conn net.Conn, port int) {
	// Defer a function to close the connection when handleConnection exits.
	defer func() {
		fmt.Println("Closing connection from port:", port)
		conn.Close() // Ensure the connection is closed
	}()

	buf := make([]byte, 1024) // Create a buffer to hold incoming data
	for {
		// Read data from the connection into the buffer
		n, err := conn.Read(buf)
		if err != nil {
			if n == 0 { // Check if the connection was closed by the client
				fmt.Println("Client closed the connection:", conn.RemoteAddr())
				return // Exit the function if the connection is closed
			}
			fmt.Println("Error reading:", err) // Log any other reading errors
			return // Exit the function on error
		}

		// Convert the received data to a string
		receivedData := string(buf[:n])
		// Create a unique transaction ID using the port and current time
		transactionID := fmt.Sprintf("tx_%d_%d", port, time.Now().UnixNano())

		// Print the received data and the port it came from
		fmt.Printf("Received data: %s, from port: %d\n", receivedData, port)

		// Send the transaction ID back to the client
		_, err = conn.Write([]byte(transactionID))
		if err != nil {
			fmt.Println("Error writing response:", err) // Log any writing errors
			return // Exit the function on error
		}
	}
}

// main function sets up the TCP server.
func main() {
	// Loop to create a TCP listener for 10 ports (8001 to 8010)
	for i := 1; i <= 10; i++ {
		port := 8000 + i // Calculate the port number (8001 to 8010)
		go func(p int) {
			// Start listening for incoming TCP connections on the specified port
			listener, err := net.Listen("tcp", ":"+strconv.Itoa(p))
			if err != nil {
				fmt.Println("Error starting server on port", p, ":", err)
				return // Exit if there is an error starting the server
			}
			defer listener.Close() // Ensure the listener is closed when the function exits

			fmt.Println("Listening on port", p) // Log the port the server is listening on

			// Infinite loop to accept incoming connections
			for {
				// Accept an incoming connection
				conn, err := listener.Accept()
				if err != nil {
					fmt.Println("Error accepting connection:", err) // Log any errors in accepting connections
					continue // Continue to the next iteration of the loop
				}
				// Start a new goroutine to handle the accepted connection
				go handleConnection(conn, p)
			}
		}(port) // Pass the port number to the goroutine
	}

	select {} // Block forever to keep the main function running
}

