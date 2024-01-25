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
	Game    *gogame.GoGame
	players map[gogame.StoneColor]Player
}

func SendResponse(conn net.Conn, response api.Response) error {
	// Create a new gob encoder with the connection as the output stream
	enc := gob.NewEncoder(conn)

	// Encode the response using the gob encoder
	err := enc.Encode(response)
	if err != nil {
		return fmt.Errorf("error encoding response: %w", err)
	}

	// write a debug message to the console
	fmt.Println("Sent response:", response)
	return nil
}

func (gh *GameHandler) Run() {
	// send a Game response to all players that are Human
	response := api.Response{Type: api.Game, Data: "Game started", Game: gh.Game}
	for _, player := range gh.players {
		if player.IsHuman() {
			err := SendResponse(player.GetConn(), response)
			if err != nil {
				fmt.Println("Error sending response:", err)
			}
		}
	}

	// start gorutines for each player so that they talk freely to game
	for _, player := range gh.players {
		go player.DoMove(*gh, player.GetConn())
	}

	// Run game
	for {
		// If next player is a AI, make a move
		if !gh.players[gh.Game.Turn].IsHuman() {
			gh.players[gh.Game.Turn].DoMove(*gh, gh.players[gh.Game.Turn].GetConn())
		}
	}
}

func (gh *GameHandler) add_player(player Player) {
	// Will be ran several times, randomize if no players in Player, else add free color to player
	if len(gh.players) == 0 {
		// Initialize Player map
		gh.players = make(map[gogame.StoneColor]Player)

		// Randomize color
		if rand.Intn(2) == 0 {
			gh.players[gogame.BLACK] = player
		} else {
			gh.players[gogame.WHITE] = player
		}
	} else {
		// if first color is taken, assign the other
		if _, ok := gh.players[gogame.BLACK]; ok {
			gh.players[gogame.WHITE] = player
		} else {
			gh.players[gogame.BLACK] = player
		}
	}
}
