// cmd/server/main.go
package main

import (
	"encoding/gob"
	"fmt"
	"gogo/internal/api"
	"gogo/pkg/gogame"
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
		}
		go handleConnection(conn)
	}
}

func decodeCommand(conn net.Conn) (api.Command, error) {
	var cmd api.Command
	err := gob.NewDecoder(conn).Decode(&cmd)
	if err != nil {
		return api.Command{}, fmt.Errorf("error decoding command: %w", err)
	}
	return cmd, nil
}

func handleConnection(conn net.Conn) {
	var gameHandler *GameHandler
	var player string
	for {
		cmd, err := decodeCommand(conn)
		if err != nil {
			fmt.Println("Failed receiving command:", err)
			break
		}
		fmt.Println("Received command:", cmd)

		// Always handle greet first to get the player name and initialize the game handler
		if gameHandler == nil && cmd.Type == api.Greet {
			player = cmd.Data
			game := gogame.NewGoGame(9, []string{player, "dumb_ai"})
			if game.MakeMove("dumb_ai", 0, 0) != nil {
				fmt.Println("Error making move for dumb_ai")
			}
			gameHandler = &GameHandler{game: game, conn: conn}
			fmt.Println("Got greeted by", cmd.Data)

			// send a OK response
			greeting := fmt.Sprintf("Welcome %s!", cmd.Data)
			response := api.Response{Type: api.Ok, Data: greeting, Game: gameHandler.game}
			err = gameHandler.send_response(response)
			if err != nil {
				fmt.Errorf("error sending response: %w", err)
			}
		} else {
			err := gameHandler.handle_cmd(cmd, player)
			if err != nil {
				fmt.Println("Error handling command:", err)
				break
			}
		}
	}
}
