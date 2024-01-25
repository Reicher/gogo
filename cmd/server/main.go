// cmd/server/main.go
package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"os"

	"github.com/Reicher/gogo/internal/api"
	"github.com/Reicher/gogo/pkg/gogame"
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
func (gh *GameHandler) AddClient(c Client) error {
	// Read the name of the client
	buf := make([]byte, 1024)
	n, err := c.conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading from client:", err)
	}
	c.name = string(buf[:n])

	// print Playername connected
	fmt.Println(c.name, "connected!")

	// TODO: Not very random
	if len(gh.players) == 0 {
		c.color = gogame.BLACK
		gh.players[gogame.BLACK] = c
	} else {
		c.color = gogame.WHITE
		gh.players[gogame.WHITE] = c
		go gh.start()
	}
	return nil
}

// Start the game
func (gh *GameHandler) start() {

	gameStartMessage := fmt.Sprintf("Starting game! %s (%s) vs. %s (%s)",
		gh.players[gogame.BLACK].name, gogame.BLACK,
		gh.players[gogame.WHITE].name, gogame.WHITE)
	fmt.Println(gameStartMessage)
	for color := range gh.players {
		player := gh.players[color]
		response := api.Response{
			Type: api.Game,
			Data: gameStartMessage,
			Game: gh.Game,
		}
		// Create a new gob encoder with the connection as the output stream
		enc := gob.NewEncoder(player.conn)

		// Encode the response using the gob encoder
		err := enc.Encode(response)
		if err != nil {
			fmt.Println("Error encoding response:", err)
		}

		go func(p Client) { p.handle(gh.Game) }(player)
	}
}

// a list of all the games handled by the server
var games []GameHandler

// Client that is waiting for a partner
var waitingClient *Client

// Handle the client and all requests from the client on the game
func (c *Client) handle(game *gogame.GoGame) {
	fmt.Println("Starting client rutine:", c.name)
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
			fmt.Println("Placing ", c.color, " stone at ", row, column)
			err = game.MakeMove(c.color, row, column)

		default:
			fmt.Println("Unknown request type:", request.Type)
		}

		// Check if the request was successful
		if err != nil {
			fmt.Println("Error making move:", err)
			response := api.Response{
				Type: api.Err,
				Data: err.Error(),
				Game: game,
			}
			// Create a new gob encoder with the connection as the output stream
			enc := gob.NewEncoder(c.conn)

			// Encode the response using the gob encoder
			err := enc.Encode(response)
			if err != nil {
				fmt.Println("Error encoding response:", err)
			}
			continue
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

// initializeGame initializes a new game, adds a client to it, and returns the created GameHandler.
func initializeGame(conn net.Conn) error {
	// Create a GoGame and a GameHandler
	game := gogame.NewGoGame(9)
	gameHandler := GameHandler{
		Game:    game,
		players: make(map[gogame.StoneColor]Client),
	}

	// Add a client to the game
	client := Client{conn: conn}
	if err := gameHandler.AddClient(client); err != nil {
		// Close the connection in case of an error
		conn.Close()
		return fmt.Errorf("could not add client to game: %w", err)
	}

	games = append(games, gameHandler)
	return nil
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
		defer conn.Close()
		if waitingClient == nil {
			// Initialize a new game and add the client to it
			waitingClient = &Client{conn: conn}
			if err := initializeGame(conn); err != nil {
				fmt.Println("Error initializing game:", err)
			}
		} else {
			// Add the client to the latest game and start a new goroutine for the game
			newClient := &Client{conn: conn}
			game := games[len(games)-1]
			if err = game.AddClient(*newClient); err != nil {
				fmt.Println("Error adding client to game:", err)
			}
			waitingClient = nil
		}
	}
}
