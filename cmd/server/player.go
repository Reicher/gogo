package main

import (
	"encoding/gob"
	"fmt"
	"gogo/internal/api"
	"math/rand"
	"net"
	"strconv"
)

type Player interface {
	GetName() string
	DoMove(game ) error
	IsHuman() bool
	GetConn() net.Conn
}

type Human struct {
	Conn net.Conn
	Name string
}

func (h *Human) decodeRequest() (api.Request, error) {
	var cmd api.Request
	err := gob.NewDecoder(h.Conn).Decode(&cmd)
	if err != nil {
		return api.Request{}, fmt.Errorf("error decoding request: %w", err)
	}
	return cmd, nil
}

func (h *Human) DoMove(game) error {
	request, err := h.decodeRequest()

	if err != nil {
		fmt.Println("Failed receiving request:", err)
	}
	if request.Type == api.Place {
		x, _ := strconv.Atoi(request.Data[:1])
		y, _ := strconv.Atoi(request.Data[2:])
		err = game.MakeMove(game.Turn, x, y)
		if err != nil {
			// send a Err response
			response := api.Response{Type: api.Err, Data: err.Error(), Game: gh.Game}
			err_response := SendResponse(h.Conn, response)
			if err_response != nil {
				return fmt.Errorf("error sending response: %w on error making move: %w", err_response, err)
			}
			return fmt.Errorf("error making move: %w", err)
		}
		// send a Ok response
		response := api.Response{Type: api.Ok, Data: "Move Ok", Game: gh.Game}
		err_response := SendResponse(h.Conn, response)
		if err_response != nil {
			return fmt.Errorf("error sending response: %w on error making move: %w", err_response, err)
		}
	}

	return nil
}

func (h *Human) GetName() string {
	return h.Name
}

func (h *Human) IsHuman() bool {
	return true
}

func (h *Human) GetConn() net.Conn {
	return h.Conn
}

type AI struct {
}

func (a *AI) DoMove(gh GameHandler, conn net.Conn) error {
	// dummy random move for dumb_ai
	x := rand.Intn(gh.Game.Board.Size)
	y := rand.Intn(gh.Game.Board.Size)
	err := gh.Game.MakeMove(gh.Game.Turn, x, y)
	if err != nil {
		return fmt.Errorf("ai error making move: %w", err)
	}
	return nil

}

func (a *AI) GetName() string {
	return "AI"
}

func (a *AI) IsHuman() bool {
	return false
}

func (a *AI) GetConn() net.Conn {
	return nil
}
