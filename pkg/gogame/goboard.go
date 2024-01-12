// internal/goboard/goboard.go
package gogame

import "fmt"

// StoneColor represents the color of a stone on the Go board.
type StoneColor string

const (
	BLACK StoneColor = "BLACK"
	WHITE StoneColor = "WHITE"
	EMPTY StoneColor = ""
)

// GoBoard represents the state of a Go board.
type GoBoard struct {
	Size  int
	Board [][]StoneColor
}

// NewGoBoard initializes a new Go board.
func NewBoard(size int) *GoBoard {
	board := make([][]StoneColor, size)
	for i := range board {
		board[i] = make([]StoneColor, size)
	}
	return &GoBoard{
		Size:  size,
		Board: board,
	}
}

// MakeMove makes a move on the Go board.
func (board *GoBoard) MakeMove(row int, column int, color StoneColor) error {
	fmt.Println(color, " Stone at row:", row, "column:", column)
	if board.isValidMove(row, column, board.Size) {
		board.Board[row][column] = color
		return nil
	}
	return fmt.Errorf("invalid move")
}

// Check if a move is valid
func (board *GoBoard) isValidMove(row, column, size int) bool {
	return (row >= 0 && row < size && column >= 0 && column < size) &&
		board.Board[row][column] == EMPTY
}

// PrintBoard prints the whole board to the terminal.
func (board *GoBoard) PrintBoard() {
	for _, row := range board.Board {
		for _, stone := range row {
			switch stone {
			case BLACK:
				fmt.Print("B ")
			case WHITE:
				fmt.Print("W ")
			case EMPTY:
				fmt.Print("- ")
			}
		}
		fmt.Println()
	}
}
