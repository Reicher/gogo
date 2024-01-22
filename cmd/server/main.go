// cmd/server/main.go
package main

import (
	"fmt"
	"gogo/internal/api"
	"gogo/pkg/gogame"
	"net"
	"os"
)

var waitingPlayers []*Human
var activeGames []*GameHandler

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

func handleConnection(conn net.Conn) {
	var gameHandler *GameHandler

	// Handle init request
	cmd, err := decodeRequest(conn)
	if err != nil {
		fmt.Println("Failed receiving request:", err)
		return
	}

	if cmd.Type != api.FindGame {
		fmt.Println("Error: First request must be of type FindGame")
		return
	}

	// extract Human player name and create game
	var game_config api.InitRequest
	err = game_config.UnmarshalJSON([]byte(cmd.Data))
	if err != nil {
		fmt.Println("Error handling init request:", err)
		os.Exit(1)
	}

	player := &Human{Conn: conn, Name: game_config.Player}

	if game_config.HumanOpponent {
		if len(waitingPlayers) > 0 {
			game := gogame.NewGoGame(game_config.Size)
			gameHandler = &GameHandler{Game: game}
			gameHandler.add_player(player)
			gameHandler.add_player(waitingPlayers[0])
			waitingPlayers = waitingPlayers[1:]

			// Add the game to the active games list
			activeGames = append(activeGames, gameHandler)

			// Start a goroutine to handle the game
			go gameHandler.Run()
		} else {
			// Add the player to the waiting players list
			waitingPlayers = append(waitingPlayers, player)
		}
	} else {
		// Create a new game with an AI player
		game := gogame.NewGoGame(game_config.Size)
		gameHandler = &GameHandler{Game: game}
		gameHandler.add_player(player)
		gameHandler.add_player(&AI{})

		// Add the game to the active games list
		activeGames = append(activeGames, gameHandler)

		// Start a goroutine to handle the game
		go gameHandler.Run()
	}
}
