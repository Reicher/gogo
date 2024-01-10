// cmd/client/main.go
package main

import (
	"encoding/gob"
	"fmt"
	"gogo/internal/goboard"
	"net"
	"os"
)

func main() {
	// Parse the two cmd line arguments.
	// server address
	// name of the player
	// If the number of arguments is not 2, print an error message and exit
	// If the server address is not valid, print an error message and exit
	// If the player name is not valid, print an error message and exit
	// Print the server address and player name
	// Connect to the server
	// Send a request to the server with the player name
	// Print the response from the server

	if len(os.Args) != 3 {
		fmt.Println("Usage: go run main.go <server address> <player name>")
		os.Exit(1)
	}

	serverAddr := os.Args[1]
	playerName := os.Args[2]

	fmt.Println("Server address:", serverAddr)
	fmt.Println("Player name:", playerName)

	// Connect to the server
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		os.Exit(1)
	} else {
		fmt.Println("Connected to server successfully!")
		defer conn.Close()
	}

	// Send a request to the server with the player name
	fmt.Fprint(conn, playerName)

	// Print the response from the server
	response := make([]byte, 1024)
	n, err := conn.Read(response)
	if err != nil {
		fmt.Println("Error reading from server:", err)
		os.Exit(1)
	}
	fmt.Println("Response from server:", string(response[:n]))

	// Create a new gob decoder with the connection as the input stream
	dec := gob.NewDecoder(conn)

	// Create a new empty goboard
	board := new(goboard.GoBoard)

	// Decode the board using the gob decoder
	err = dec.Decode(board)
	if err != nil {
		fmt.Println("Error decoding the board:", err)
	} else {
		fmt.Println("Board decoded successfully!")
		board.PrintBoard()
	}
}
