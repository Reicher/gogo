// cmd/server/main.go
package main

import (
	"encoding/gob"
	"fmt"
	"gogo/internal/api"
	"gogo/pkg/goboard"
	"math/rand"
	"net"
	"os"
	"strconv"
)

func decodeCommand(conn net.Conn) (api.Command, error) {
	var cmd api.Command
	err := gob.NewDecoder(conn).Decode(&cmd)
	if err != nil {
		return api.Command{}, fmt.Errorf("error decoding command: %w", err)
	}
	return cmd, nil
}

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
		}
		go handleConnection(conn, goboard.NewBoard(size))
	}
}

func handleConnection(conn net.Conn, board *goboard.GoBoard) {
	fmt.Println("Someone connected and want to say something!")

	for {
		cmd, err := decodeCommand(conn)
		if err != nil {
			fmt.Println("Failed reciving command:", err)
			continue
		}

		fmt.Println("Received command:", cmd)

		// Handle the command based on its type
		switch cmd.Type {
		case api.Greet:
			fmt.Println("Got greeted!")
		case api.GetBoard:
			// Create a new gob encoder with the connection as the output stream
			enc := gob.NewEncoder(conn)

			// Encode the board using the gob encoder
			err := enc.Encode(board)
			if err != nil {
				fmt.Println("Error encoding the board:", err)
			} else {
				fmt.Println("Board encoded successfully!")
			}
		case api.MakeMove:
			// split data at " " to get x and y
			x, _ := strconv.Atoi(cmd.Data[:1])
			y, _ := strconv.Atoi(cmd.Data[2:])
			board.MakeMove(x, y, goboard.WHITE)
			board.MakeMove(rand.Intn(board.Size), rand.Intn(board.Size), goboard.BLACK)

		case api.Pass:
			// Handle the Pass command
		case api.Resign:
			// Handle the Resign command
		default:
			fmt.Println("Received unknown command type:", cmd.Type)
		}

		// Make a random move between 0 and board size
		board.PrintBoard()
	}
}
