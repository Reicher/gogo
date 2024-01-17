package api

import (
	"gogo/pkg/gogame"
)

type ResponseType int

const (
	Ok ResponseType = iota
	Err
	Game
)

type Response struct {
	Type ResponseType
	Data string
	Game *gogame.GoGame
}
