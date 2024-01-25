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

func send_cmd(conn net.Conn, cmd api.Request) error {
	fmt.Println("Sending request:", cmd)

	// Create a gob encoder and send the request
	enc := gob.NewEncoder(conn)
	err := enc.Encode(&cmd)
	if err != nil {
		return fmt.Errorf("error encoding request: %w", err)
	}
	return nil
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run main.go <server address> <player name>")
		os.Exit(1)
	}

	// Initialize the board
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

	// Send the player name to the server
	_, err = conn.Write([]byte(playerName))
	if err != nil {
		fmt.Println("Error sending player name:", err)
		os.Exit(1)
	}

	// Wait for the server to send a response
	var response *api.Response
	dec := gob.NewDecoder(conn)
	err = dec.Decode(&response)
	if err != nil {
		fmt.Println("Error decoding response:", err)
		os.Exit(1)
	}
	response.Game.Print()

	go handleGameResponse(conn)

	// Send requests to the server
	for {
		cmd := get_cmd_from_console()
		err = send_cmd(conn, cmd)
		if err != nil {
			fmt.Println("Error sending request:", err)
			os.Exit(1)
		}
	}
}

func handleGameResponse(conn net.Conn) {
	for {
		dec := gob.NewDecoder(conn)
		var response api.Response
		err := dec.Decode(&response)

		fmt.Print("\033[H\033[2J") // clear the screen

		if err != nil {
			fmt.Println("Error decoding response:", err)
			os.Exit(1)
		}
		// Check if the response is an error
		if response.Type == api.Err {
			fmt.Println("Error:", response.Data)
			continue
		}
		response.Game.Print()
		fmt.Println("Available requests:")
		fmt.Println("  Place Stone <row> <column>")
		fmt.Println("  Pass")
		fmt.Println("  Resign")
		fmt.Print("Enter request: ")
	}
}

func get_cmd_from_console() api.Request {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("Available requests:")
		fmt.Println("  Place Stone <row> <column>")
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
