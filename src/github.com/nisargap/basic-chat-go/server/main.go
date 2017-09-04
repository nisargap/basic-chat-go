package main

// This is a basic implementation of a server-side Go websocket server
// References: https://scotch.io/bar-talk/build-a-realtime-chat-server-with-go-and-websockets
import (
	"github.com/gorilla/websocket"
	"github.com/rs/cors"
	"log"
	"net/http"
)

// Global variables are usually a bad idea but for learning purposes it should be fine
// Stores the websocket connection as a key and the value being a bool
var clients = make(map[*websocket.Conn]bool)

// Acts as a queue for messages sent
var broadcast = make(chan Message)

// Create an instance of a websocket upgrader to upgrade the HTTP request to
// websocket
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	}}

// We only need UserId and Message on the backend
type Message struct {
	UserId  string `json:"username"`
	Message string `json:"message"`
}

// This function extracts the message from the broadcast channel global
// variable that was created and then loops over the clients global map
// and takes the connection key and sends out the message to ever client
// that is currently connected
func handleMessages() {
	// An infinite loop is needed here because the goroutine is constantly
	// reading things from the broadcast variable
	for {
		// Grab the next message from the broadcast channel
		msg := <-broadcast
		//Log the message to the console
		log.Println(msg)
		// Send it out to every client that is currently connected
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	// Make sure we close the connection when the function returns
	defer ws.Close()
	// add the ws to the clients global var - very important
	clients[ws] = true
	for {
		var msg Message
		// Read in a new message as JSON and map it to a Message object
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, ws)
			break
		}
		// Send the newly received message to the broadcast channel
		broadcast <- msg
	}
}

// Entrypoint function
func main() {
	// Creating a goroutine that runs alongside the server
	// This goroutine listens for incoming messages on the channel
	go handleMessages()
	mux := http.NewServeMux()
	// Handle websocket requests
	mux.HandleFunc("/v1/ws", handleConnections)
	// Friendly message for running server
	log.Println("Listening on port localhost:3030...")
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
	})

	handler := c.Handler(mux)
	// Run server on port 3030
	err := http.ListenAndServe(":3030", handler)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
