// cmd/server/main.go
package main

import (
	"fmt"
	"gogo/internal/api"
	"gogo/pkg/gogame"
	"net"
	"os"
)

func main() {
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

func handle_init_request(init_data string, conn net.Conn) (*GameHandler, error) {
	var initData api.InitRequest
	err := initData.UnmarshalJSON([]byte(init_data))
	if err != nil {
		fmt.Println("Error unmarshalling request:", err)
	}

	// extract name and create Human player
	player1 := Human{Conn: conn, Name: initData.Player}
	game := gogame.NewGoGame(initData.Size)

	// create a ai-game
	if initData.HumanOpponent {
		return nil, fmt.Errorf("Human opponent not implemented yet")
	}

	// create game handler
	gamehandler := GameHandler{Game: game, conn: conn}
	gamehandler.add_player(&player1)
	gamehandler.add_player(&AI{})

	// send response
	greeting := fmt.Sprintf("Welcome %s!", initData.Player)
	response := api.Response{Type: api.Ok, Data: greeting, Game: game}
	err = gamehandler.send_response(response)
	if err != nil {
		return nil, fmt.Errorf("error sending response: %w", err)
	}

	return &gamehandler, nil
}

func handleConnection(conn net.Conn) {
	var gameHandler *GameHandler

	// Handle init request
	cmd, err := decodeRequest(conn)
	if err != nil {
		fmt.Println("Failed receiving request:", err)
		os.Exit(1)
	}
	if cmd.Type == api.FindGame {
		gameHandler, err = handle_init_request(cmd.Data, conn)
		if err != nil {
			fmt.Println("Error handling init request:", err)
			os.Exit(1)
		}
	} else {
		fmt.Println("Error: First request must be of type FindGame")
		os.Exit(1)
	}

	// Main game loop
	for {
		next_player := gameHandler.Player[gameHandler.Game.Turn]
		err = next_player.DoMove(*gameHandler, conn)
		if err != nil {
			fmt.Println("Error making move:", err)
			os.Exit(1)
		}
	}
}
