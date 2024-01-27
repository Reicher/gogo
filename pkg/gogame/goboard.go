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
	Size       int
	Board      [][]StoneColor
	last_board [][]StoneColor
}

// copyBoard copies the board to a new board using copy(to, from)
func (board *GoBoard) copyBoard() *[][]StoneColor {
	current_board := board.Board
	new_board := make([][]StoneColor, len(current_board))
	for i := range current_board {
		new_board[i] = make([]StoneColor, len(current_board[i]))
		copy(new_board[i], current_board[i])
	}

	return &new_board
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
func (board *GoBoard) PutStone(row int, column int, color StoneColor) error {
	fmt.Println(color, " Stone at row:", row, "column:", column)
	if board.isValidMove(row, column, color) {
		// create a new copied board that looks like the old
		board.last_board = *board.copyBoard()
		board.Board[row][column] = color
		return nil
	}
	return fmt.Errorf("invalid move")
}

// Check if a move is valid
func (board *GoBoard) isValidMove(row int, column int, color StoneColor) bool {
	unused := (row >= 0 && row < board.Size && column >= 0 && column < board.Size) &&
		board.Board[row][column] == EMPTY

	return unused && board.ko(row, column, color)
}

// Check if a move is a ko
func (board *GoBoard) ko(row, column int, color StoneColor) bool {
	new_board := *board.copyBoard()
	new_board[row][column] = color
	// check if the new move would create an identical board to the last one
	for i := range board.Board {
		for j := range board.Board[i] {
			if new_board[i][j] != board.last_board[i][j] {
				return true
			}
		}
	}
	return true
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
