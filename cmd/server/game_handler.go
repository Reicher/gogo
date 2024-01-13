// cmd/server/gamehandler.go
package main

import (
	"encoding/gob"
	"fmt"
	"gogo/internal/api"
	"gogo/pkg/gogame"
	"math/rand"
	"net"
	"strconv"
)

type GameHandler struct {
	game *gogame.GoGame
	conn net.Conn
}

func (gh *GameHandler) send_response(response api.Response) error {
	// Create a new gob encoder with the connection as the output stream
	enc := gob.NewEncoder(gh.conn)

	// Encode the response using the gob encoder
	err := enc.Encode(response)
	if err != nil {
		return fmt.Errorf("error encoding response: %w", err)
	} else {
		fmt.Println("Response encoded successfully!")
	}
	return nil
}

func (gh *GameHandler) handleGreet() {
	fmt.Println("Got greeted again!?!")
}

func (gh *GameHandler) handleMakeMove(cmd api.Command, player string) error {
	// split data at " " to get x and y
	x, _ := strconv.Atoi(cmd.Data[:1])
	y, _ := strconv.Atoi(cmd.Data[2:])
	return gh.game.MakeMove(player, x, y)
}

func (gh *GameHandler) handlePass() {
	// Handle the Pass command
}

func (gh *GameHandler) handle_cmd(cmd api.Command, player string) error {
	// Handle the command based on its type
	var err error = nil
	var response api.Response

	switch cmd.Type {
	case api.Greet:
		gh.handleGreet()
		// send a OK response
		greeting := fmt.Sprintf("Welcome %s!", cmd.Data)
		response = api.Response{Type: api.Ok, Data: greeting, Game: gh.game}
		err = gh.send_response(response)
		if err != nil {
			return fmt.Errorf("error sending response: %w", err)
		}
	case api.MakeMove:
		err = gh.handleMakeMove(cmd, player)
		if err != nil {
			// send a Err response
			response = api.Response{Type: api.Err, Data: err.Error(), Game: gh.game}
			err = gh.send_response(response)
			if err != nil {
				return fmt.Errorf("error sending response: %w", err)
			}
			return fmt.Errorf("error making move: %w", err)
		}

		// dummy random move for dumb_ai
		x := rand.Intn(gh.game.Board.Size)
		y := rand.Intn(gh.game.Board.Size)
		cmd := api.Command{Type: api.MakeMove, Data: fmt.Sprintf("%d %d", x, y)}
		gh.handleMakeMove(cmd, "dumb_ai")

		// send a OK response
		response = api.Response{Type: api.Ok, Data: "Move made successfully!", Game: gh.game}
		err = gh.send_response(response)
	case api.Pass:
		gh.handlePass()
	}
	return err
}
