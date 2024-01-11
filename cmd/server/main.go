// cmd/server/main.go
package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"gogo/internal/api"
	"gogo/pkg/goboard"
	"math/rand"
	"net"
	"os"
	"strconv"
)

func decodeCommand(conn net.Conn) api.Command {
	data := make([]byte, 1024)
	_, err := conn.Read(data)
	if err != nil {
		fmt.Println("Error reading data from client:", err)
		os.Exit(1)
	}

	var cmd api.Command
	dec := gob.NewDecoder(bytes.NewReader(data))
	err = dec.Decode(&cmd)
	if err != nil {
		fmt.Println("Error decoding command:", err)
		os.Exit(1)
	}
	return cmd
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
		cmd := decodeCommand(conn)
		fmt.Println("Received command:", cmd)

		// Handle the command based on its type
		switch cmd.Type {
		case api.Greet:
			fmt.Println("Got greeted!")
			break
		case api.GetBoard:
			// Create a new gob encoder with the connection as the output stream
			enc := gob.NewEncoder(conn)

			// Encode the board using the gob encoder
			err := enc.Encode(board)
			if err != nil {
				fmt.Println("Error encoding the board:", err)
			} else {
				fmt.Println("Board encoded successfully!")
				board.PrintBoard()
			}
			break
		case api.MakeMove:
			x := int(cmd.Data[0])
			y := int(cmd.Data[1])
			board.MakeMove(x, y, goboard.WHITE)
		case api.Pass:
			// Handle the Pass command
		case api.Resign:
			// Handle the Resign command
			break
		default:
			fmt.Println("Received unknown command type:", cmd.Type)
		}

		// Make a random move between 0 and board size
		board.MakeMove(rand.Intn(board.Size), rand.Intn(board.Size), goboard.BLACK)
	}
}
