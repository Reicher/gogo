// cmd/server/main.go
package main

import (
	"encoding/gob"
	"fmt"
	"gogo/internal/goboard"
	"net"
	"os"
	"strconv"
)

func main() {
	// Default board size
	size := 9

	// Check if a command-line argument for size is provided
	if len(os.Args) > 1 {
		size, _ = strconv.Atoi(os.Args[1])
	}
	fmt.Println("Using board size:", size)

	// Start listening for connections
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error starting the server:", err)
		os.Exit(1)
	}
	fmt.Println("Server started successfully!")

	// For each connection, create a new goroutine
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			os.Exit(1)
		}
		go handleConnection(conn, goboard.NewBoard(size))
	}
}

func handleConnection(conn net.Conn, board *goboard.GoBoard) {
	// In the goroutine, read the player name from the connection
	playerName := make([]byte, 1024)
	n, err := conn.Read(playerName)
	if err != nil {
		fmt.Println("Error reading player name:", err)
		os.Exit(1)
	}
	// Print the player name and the welcome message
	fmt.Println("Player", string(playerName[:n]), "connected!")

	welcomeMessage := fmt.Sprintf("Hello %s! Welcome to the game!", string(playerName[:n]))
	fmt.Fprint(conn, welcomeMessage)
	fmt.Println("Sent welcome message:", welcomeMessage)

	// Make a move for fun and see that it comes through on the client side
	board.MakeMove(4, 4, goboard.BLACK)

	// Create a new gob encoder with the connection as the output stream
	enc := gob.NewEncoder(conn)

	// Encode the board using the gob encoder
	err = enc.Encode(board)
	if err != nil {
		fmt.Println("Error encoding the board:", err)
	} else {
		fmt.Println("Board encoded successfully!")
		board.PrintBoard()
	}
}
