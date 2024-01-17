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
	DoMove(gh GameHandler, conn net.Conn) error
}

type Human struct {
	Conn net.Conn
	Name string
}

func decodeRequest(conn net.Conn) (api.Request, error) {
	var cmd api.Request
	err := gob.NewDecoder(conn).Decode(&cmd)
	if err != nil {
		return api.Request{}, fmt.Errorf("error decoding request: %w", err)
	}
	return cmd, nil
}

func (h *Human) DoMove(gh GameHandler, conn net.Conn) error {
	request, err := decodeRequest(conn)

	if err != nil {
		fmt.Println("Failed receiving request:", err)
	}
	if request.Type == api.Place {
		x, _ := strconv.Atoi(request.Data[:1])
		y, _ := strconv.Atoi(request.Data[2:])
		err = gh.Game.MakeMove(gh.Game.Turn, x, y)
		if err != nil {
			// send a Err response
			response := api.Response{Type: api.Err, Data: err.Error(), Game: gh.Game}
			err_response := gh.send_response(response)
			if err_response != nil {
				return fmt.Errorf("error sending response: %w on error making move: %w", err_response, err)
			}
			return fmt.Errorf("error making move: %w", err)
		}
		// send a Ok response
		response := api.Response{Type: api.Ok, Data: "Move Ok", Game: gh.Game}
		err_response := gh.send_response(response)
		if err_response != nil {
			return fmt.Errorf("error sending response: %w on error making move: %w", err_response, err)
		}
	}

	return nil
}

func (h *Human) GetName() string {
	return h.Name
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
