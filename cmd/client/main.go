// cmd/client/main.go
package main

import (
	"bufio"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"gogo/internal/api"
	"net"
	"os"
	"strings"
	"time"
)

func send_cmd(conn net.Conn, cmd api.Request) (*api.Response, error) {
	fmt.Println("Sending request:", cmd)

	// Create a gob encoder and send the request
	enc := gob.NewEncoder(conn)
	err := enc.Encode(&cmd)
	if err != nil {
		return nil, fmt.Errorf("error encoding request: %w", err)
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

// create a json-string on the format api.InitRequest (size, player, human_opponent)
func initialize_game(conn net.Conn, size int, player_name string, human_opponent bool) (response *api.Response, err error) {
	fmt.Println("Initializing game:", size, player_name, human_opponent)

	// Create an api.InitRequest object
	initReq := api.InitRequest{
		Size:          size,
		Player:        player_name,
		HumanOpponent: human_opponent,
	}

	// Serialize the api.InitRequest object to JSON
	jsonInit, err := json.Marshal(initReq)
	if err != nil {
		return nil, fmt.Errorf("error marshalling request: %w", err)
	}

	// Send the request
	response, err = send_cmd(conn, api.Request{Type: api.FindGame, Data: string(jsonInit)})
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}

	// Print the response
	fmt.Println("Response:", response.Data)
	return response, err

}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run main.go <server address> <player name>")
		os.Exit(1)
	}

	// Initialize the board
	serverAddr := os.Args[1]
	playerName := os.Args[2]
	human_opponent := true

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

	// Send request to initialize the game
	response, err = initialize_game(conn, 9, playerName, human_opponent)
	if err != nil {
		fmt.Println("Error initializing game:", err)
		os.Exit(1)
	}

	for {
		fmt.Print("\033[H\033[2J") // clear the screen
		if response.Type != api.Ok {
			fmt.Println("Error:", response.Data)
		}
		response.Game.Print()

		cmd := get_cmd_from_console()
		response, err = send_cmd(conn, cmd)

		fmt.Println("Response:", response.Data)

		if err != nil {
			fmt.Println("Error sending request:", err)
			os.Exit(1)
		}

		// TODO: Add some kind of retry logic here
		// sleep for a second
		time.Sleep(1 * time.Second)

		// send after game update
		response, err = send_cmd(conn, api.Request{Type: api.Update})
		if err != nil {
			fmt.Println("Error sending request:", err)
			os.Exit(1)
		}

		response.Game.Print()
	}
}

func get_cmd_from_console() api.Request {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("Available requests:")
		fmt.Println("  MakeMove <row> <column>")
		fmt.Println("  Pass")
		fmt.Println("  Resign")
		fmt.Print("Enter request: ")

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		parts := strings.SplitN(input, " ", 2)
		if len(parts) < 2 {
			fmt.Println("Invalid request, please enter a request in the format 'type data'")
			continue
		}

		requestType, ok := api.StringToRequestType(parts[0])
		if !ok {
			fmt.Println("Unknown request:", parts[0])
			continue
		}

		return api.Request{
			Type: requestType,
			Data: parts[1],
		}
	}
}
