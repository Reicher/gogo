// cmd/client/main.go
package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"gogo/internal/api"
	"gogo/pkg/goboard"
	"net"
	"os"
	"strings"
)

var board *goboard.GoBoard

func send_cmd(conn net.Conn, cmd api.Command) error {
	fmt.Println("Sending command:", cmd)

	// Create a gob encoder and send the command
	enc := gob.NewEncoder(conn)
	err := enc.Encode(&cmd)
	if err != nil {
		return fmt.Errorf("error encoding command: %w", err)
	}

	return nil
}

// Function to receive board updates from the server
func update_board(conn net.Conn) {
	dec := gob.NewDecoder(conn)
	err := dec.Decode(board)
	if err != nil {
		fmt.Println("Error decoding the board:", err)
	} else {
		fmt.Println("Board decoded successfully!")
		board.PrintBoard()
	}
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run main.go <server address> <player name>")
		os.Exit(1)
	}

	// Initialize the board
	board = &goboard.GoBoard{}

	serverAddr := os.Args[1]
	playerName := os.Args[2]

	// Connect to the server
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		os.Exit(1)
	} else {
		fmt.Println("Connected to server successfully!")
		defer conn.Close()
	}

	err = send_cmd(conn, api.Command{Type: api.Greet, Data: playerName})
	if err != nil {
		fmt.Println("Error sending command:", err)
		os.Exit(1)
	}

	// Get initial board
	send_cmd(conn, api.Command{Type: api.GetBoard, Data: ""})
	update_board(conn)

	for {
		cmd := get_cmd_from_console()
		send_cmd(conn, cmd)

		send_cmd(conn, api.Command{Type: api.GetBoard, Data: ""})
		update_board(conn)
	}
}

func get_cmd_from_console() api.Command {
	reader := bufio.NewReader(os.Stdin)
	for {
		// presents all avaiable commands to the user
		// Present a list of commands to the user
		// User is then prompted for any data required by the command
		// The command and data are returned

		fmt.Println("Available commands:")
		fmt.Println("  Greet <player name>")
		fmt.Println("  GetBoard")
		fmt.Println("  MakeMove <row> <column>")
		fmt.Println("  pass")
		fmt.Println("  Resign")
		fmt.Print("Enter command: ")

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		parts := strings.SplitN(input, " ", 2)
		if len(parts) < 2 {
			fmt.Println("Invalid command, please enter a command in the format 'command data'")
			continue
		}

		commandType, ok := api.StringToCommandType(parts[0])
		if !ok {
			fmt.Println("Unknown command:", parts[0])
			continue
		}

		return api.Command{
			Type: commandType,
			Data: parts[1],
		}
	}
}
