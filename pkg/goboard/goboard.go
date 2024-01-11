// internal/goboard/goboard.go
package goboard

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
func (board *GoBoard) MakeMove(row int, column int, color StoneColor) {
	fmt.Println(color, " Stone at row:", row, "column:", column)
	if isValidMove(row, column, board.Size) {
		board.Board[row][column] = color
	}
}

// Score calculates and returns the score of the Go board.
func (board *GoBoard) Score() int {
	// Implement the scoring logic here
	return 0
}

// Helper function to check if a move is valid
func isValidMove(row, column, size int) bool {
	return row >= 0 && row < size && column >= 0 && column < size
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
