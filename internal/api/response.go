package api

import (
	"github.com/Reicher/gogo/pkg/gogame"
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

// Marshal a response to a string
func (r *Response) Marshal() string {
	return r.Data
}
