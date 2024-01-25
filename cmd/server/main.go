// cmd/server/main.go
package main

import (
	"encoding/gob"
	"fmt"
	"gogo/internal/api"
	"gogo/pkg/gogame"
	"net"
	"os"
)

// Struct of a client
type Client struct {
	conn  net.Conn
	name  string
	color gogame.StoneColor
}

type GameHandler struct {
	Game    *gogame.GoGame
	players map[gogame.StoneColor]Client
}

// Add a client to the game, assigning a random color on the first client and the opposite color on the second client
func (gh *GameHandler) AddClient(c *Client) {
	// Read the name of the client
	buf := make([]byte, 1024)
	n, err := c.conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading from client:", err)
		return
	}
	c.name = string(buf[:n])

	fmt.Println("Adding client:", c.name)

	// TODO: Not very random
	if len(gh.players) == 0 {
		gh.players[gogame.BLACK] = *c
		c.color = gogame.BLACK
	} else {
		gh.players[gogame.WHITE] = *c
		c.color = gogame.WHITE
		go gh.start()
	}
}

// Start the game
func (gh *GameHandler) start() {
	fmt.Println("Starting game!")

	// send the board to both players
	for _, player := range gh.players {
		response := api.Response{
			Type: api.Game,
			Data: "Game started!",
			Game: gh.Game,
		}
		// Create a new gob encoder with the connection as the output stream
		enc := gob.NewEncoder(player.conn)

		// Encode the response using the gob encoder
		err := enc.Encode(response)
		if err != nil {
			fmt.Println("Error encoding response:", err)
		}
		fmt.Println("Sending response:", response)

		// Start a gorutine for the player to handle requests
		go player.handle(gh.Game)
	}
}

// a list of all the games handled by the server
var games []GameHandler

// Client that is waiting for a partner
var waitingClient *Client

// Handle the client and all requests from the client on the game
func (c *Client) handle(game *gogame.GoGame) {
	for {
		// wait for a api.request from the client
		request := api.Request{}
		dec := gob.NewDecoder(c.conn)
		err := dec.Decode(&request)
		if err != nil {
			fmt.Println("Error decoding request:", err)
			return
		}

		// Do the request on the game
		switch request.Type {
		case api.Place:
			// split the data into row and column
			var row, column int
			fmt.Sscanf(request.Data, "%d %d", &row, &column)
			game.MakeMove(c.color, row, column)
		default:
			fmt.Println("Unknown request type:", request.Type)
		}

		// Send the updated game to both players
		for _, player := range games[0].players {
			response := api.Response{
				Type: api.Game,
				Data: "Game updated!",
				Game: game,
			}
			// Create a new gob encoder with the connection as the output stream
			enc := gob.NewEncoder(player.conn)

			// Encode the response using the gob encoder
			err := enc.Encode(response)
			if err != nil {
				fmt.Println("Error encoding response:", err)
			}
			fmt.Println("Sending response:", response)
		}
	}
}

func main() {
	// Start listening for connections
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error starting the server:", err)
		os.Exit(1)
	}
	fmt.Println("Server started successfully!")

	// If a client connects and waitingClient is nil, set waitingClient to the client
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
		}
		if waitingClient == nil {
			// Create a go game and add it to games and add the client to the game
			waitingClient = &Client{conn: conn}
			game := gogame.NewGoGame(9)
			game_handler := GameHandler{Game: game, players: make(map[gogame.StoneColor]Client)}
			game_handler.AddClient(waitingClient)
			games = append(games, game_handler)
		} else {
			// Add the client to the latest game and start start a new gorutine for the game
			new_client := &Client{conn: conn}
			game := games[len(games)-1]
			game.AddClient(new_client)
			waitingClient = nil
		}
	}
}
