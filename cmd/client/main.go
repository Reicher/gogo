// cmd/client/main.go
package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"gogo/internal/api"
	"net"
	"os"
	"strings"
)

func send_cmd(conn net.Conn, cmd api.Command) (*api.Response, error) {
	fmt.Println("Sending command:", cmd)

	// Create a gob encoder and send the command
	enc := gob.NewEncoder(conn)
	err := enc.Encode(&cmd)
	if err != nil {
		return nil, fmt.Errorf("error encoding command: %w", err)
	}

	// Create a gob decoder and decode the response
	dec := gob.NewDecoder(conn)
	var response api.Response
	err = dec.Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &response, nil
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run main.go <server address> <player name>")
		os.Exit(1)
	}

	// Initialize the board
	serverAddr := os.Args[1]
	playerName := os.Args[2]

	var response *api.Response

	// Connect to the server
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		os.Exit(1)
	} else {
		fmt.Println("Connected to server successfully!")
		defer conn.Close()
	}

	// Send the greet command to initialize the game
	response, err = send_cmd(conn, api.Command{Type: api.Greet, Data: playerName})
	if err != nil {
		fmt.Println("Error sending command:", err)
		os.Exit(1)
	}

	// Print the response
	fmt.Println("Response:", response.Data)
	response.Game.Print()

	for {
		cmd := get_cmd_from_console()

		// clear the screen
		fmt.Print("\033[H\033[2J")
		response, err = send_cmd(conn, cmd)
		fmt.Println("Response:", response.Data)

		if err != nil {
			fmt.Println("Error sending command:", err)
			os.Exit(1)
		}
		response.Game.Print()
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
		fmt.Println("  MakeMove <row> <column>")
		fmt.Println("  Pass")
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
