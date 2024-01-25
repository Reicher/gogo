// internal/goboard/goboard.go
package gogame

import (
	"fmt"
)

// GoGame represents the state of a Go game.
// a board, a turn and a map of stone color to player
type GoGame struct {
	Board GoBoard
	Turn  StoneColor
}

// NewGoGame initializes a new Go game, with two player names
func NewGoGame(size int) *GoGame {

	goGame := GoGame{
		Board: *NewBoard(size),
		Turn:  BLACK,
	}

	return &goGame
}

// MakeMove makes a move on the Go board.
func (game *GoGame) MakeMove(color StoneColor, row int, column int) error {
	// Check if it is the players turn
	if game.Turn != color {
		return fmt.Errorf("not your turn")
	}

	// Make move
	err := game.Board.PutStone(row, column, game.Turn)
	if err != nil {
		return fmt.Errorf("could not put stone: %w", err)
	}

	// Switch turn
	if game.Turn == BLACK {
		game.Turn = WHITE
	} else {
		game.Turn = BLACK
	}
	return nil
}

// print the state of the game
func (game *GoGame) Print() {
	game.Board.PrintBoard()
	fmt.Println("\nCurrent turn:", game.Turn)
}

func (game *GoGame) GetBoard() GoBoard {
	return game.Board
}
