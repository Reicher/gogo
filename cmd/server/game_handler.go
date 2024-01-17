// cmd/server/gamehandler.go
package main

import (
	"encoding/gob"
	"fmt"
	"gogo/internal/api"
	"gogo/pkg/gogame"
	"math/rand"
	"net"
)

type GameHandler struct {
	Game   *gogame.GoGame
	Player map[gogame.StoneColor]Player
	conn   net.Conn
}

func (gh *GameHandler) add_player(player Player) {
	// Will be ran several times, randomize if no players in Player, else add free color to player
	if len(gh.Player) == 0 {
		// Initialize Player map
		gh.Player = make(map[gogame.StoneColor]Player)

		// Randomize color
		if rand.Intn(2) == 0 {
			gh.Player[gogame.BLACK] = player
		} else {
			gh.Player[gogame.WHITE] = player
		}
	} else {
		// if first color is taken, assign the other
		if _, ok := gh.Player[gogame.BLACK]; ok {
			gh.Player[gogame.WHITE] = player
		} else {
			gh.Player[gogame.BLACK] = player
		}
	}
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
