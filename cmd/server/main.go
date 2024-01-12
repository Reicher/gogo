// cmd/server/main.go
package main

import (
	"encoding/gob"
	"fmt"
	"gogo/internal/api"
	"gogo/pkg/gogame"
	"math/rand"
	"net"
	"os"
	"strconv"
)

type GameHandler struct {
	game *gogame.GoGame
	conn net.Conn
}

func (gh *GameHandler) handleGreet() {
	fmt.Println("Got greeted again!?!")
}

func (gh *GameHandler) handleGetGame() error {
	// Create a new gob encoder with the connection as the output stream
	enc := gob.NewEncoder(gh.conn)

	// Encode the board using the gob encoder
	err := enc.Encode(gh.game)
	if err != nil {
		return fmt.Errorf("error encoding board: %w", err)
	} else {
		fmt.Println("Board encoded successfully!")
	}
	return nil
}

func (gh *GameHandler) handleMakeMove(cmd api.Command, player string) {
	// split data at " " to get x and y
	x, _ := strconv.Atoi(cmd.Data[:1])
	y, _ := strconv.Atoi(cmd.Data[2:])
	gh.game.MakeMove(player, x, y)
}

func (gh *GameHandler) handlePass() {
	// Handle the Pass command
}

func (gh *GameHandler) handle_cmd(cmd api.Command, player string) error {
	// Handle the command based on its type
	switch cmd.Type {
	case api.Greet:
		gh.handleGreet()
	case api.GetGame:
		return gh.handleGetGame()
	case api.MakeMove:
		gh.handleMakeMove(cmd, player)
		// dummy random move for dumb_ai
		x := rand.Intn(gh.game.Board.Size)
		y := rand.Intn(gh.game.Board.Size)
		cmd := api.Command{Type: api.MakeMove, Data: fmt.Sprintf("%d %d", x, y)}
		gh.handleMakeMove(cmd, "dumb_ai")
	case api.Pass:
		gh.handlePass()
	}
	return nil
}

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
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	fmt.Println("Someone connected and want to say something!")

	var gameHandler *GameHandler
	var player string
	for {
		cmd, err := decodeCommand(conn)
		if err != nil {
			fmt.Println("Failed receiving command:", err)
			break
		}

		fmt.Println("Received command:", cmd)

		if gameHandler == nil && cmd.Type == api.Greet {
			player = cmd.Data
			game := gogame.NewGoGame(9, []string{player, "dumb_ai"})
			if game.MakeMove("dumb_ai", 0, 0) != nil {
				fmt.Println("Error making move for dumb_ai")
			}
			gameHandler = &GameHandler{game: game, conn: conn}
			fmt.Println("Got greeted by", cmd.Data)
		} else {
			err := gameHandler.handle_cmd(cmd, player)
			if err != nil {
				fmt.Println("Error handling command:", err)
				break
			}
		}

		if gameHandler.game != nil {
			gameHandler.game.Print()
		}
	}
}
