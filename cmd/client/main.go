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

// Sends a command to the server and receives the updated board
func send_cmd(conn net.Conn, cmd api.Command) {
	fmt.Println("Sending command:", cmd)

	// Create a gob encoder and send the command
	enc := gob.NewEncoder(conn)
	err := enc.Encode(&cmd)
	if err != nil {
		fmt.Println("Error encoding command:", err)
		os.Exit(1)
	}

	// Create a new gob decoder with the connection as the input stream
	dec := gob.NewDecoder(conn)

	// Decode the board using the gob decoder
	err = dec.Decode(board)

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

	serverAddr := os.Args[1]
	playerName := os.Args[2]

	fmt.Println("Server address:", serverAddr)
	fmt.Println("Player name:", playerName)

	// Connect to the server
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		os.Exit(1)
	} else {
		fmt.Println("Connected to server successfully!")
		defer conn.Close()
	}

	send_cmd(conn, api.Command{Type: api.Greet, Data: playerName})

	reader := bufio.NewReader(os.Stdin)
	for {
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

		cmd := api.Command{
			Type: commandType,
			Data: parts[1],
		}
		send_cmd(conn, cmd)

		send_cmd(conn, api.Command{Type: api.GetBoard})
		board.PrintBoard()
	}
}
