// internal/goboard/goboard.go
package gogame

import (
	"fmt"
	"math/rand"
)

// GoGame represents the state of a Go game.
// a board, a turn and a map of stone color to player
type GoGame struct {
	Board  GoBoard
	Turn   StoneColor
	Player map[StoneColor]string
}

// NewGoGame initializes a new Go game, with two player names
func NewGoGame(size int, player []string) *GoGame {

	goGame := GoGame{
		Board:  *NewBoard(size),
		Turn:   BLACK,
		Player: make(map[StoneColor]string),
	}

	// Randomize player stone color
	if rand.Intn(2) == 0 {
		goGame.Player[BLACK] = player[0]
		goGame.Player[WHITE] = player[1]
	} else {
		goGame.Player[BLACK] = player[1]
		goGame.Player[WHITE] = player[0]
	}

	return &goGame
}

// MakeMove makes a move on the Go board.
func (game *GoGame) MakeMove(player string, row int, column int) error {
	// Check if it is the players turn
	if game.Player[game.Turn] != player {
		return fmt.Errorf("not your turn")
	}

	// Make move
	err := game.Board.MakeMove(row, column, game.Turn)
	if err != nil {
		return err
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
	fmt.Println("Current turn:", game.Turn)
}

func (game *GoGame) GetBoard() GoBoard {
	return game.Board
}
