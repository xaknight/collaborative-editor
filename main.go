package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// We'll use a simple string to represent the shared document content.
var documentContent = ""

// Keep track of connected clients.
var clients = make(map[*websocket.Conn]bool)

// Channel to broadcast updates to all clients.
var broadcast = make(chan string)

// Configure the WebSocket upgrader.
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow connections from any origin (for development only!)
	},
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to a WebSocket connection.
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	// Register the new client.
	clients[ws] = true

	// Send the current document content to the newly connected client.
	if err := ws.WriteMessage(websocket.TextMessage, []byte(documentContent)); err != nil {
		log.Println(err)
	}

	for {
		// Read messages from the client.
		_, msg, err := ws.ReadMessage()
		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, ws)
			break
		}

		// Update the document content (very basic conflict resolution - last write wins).
		documentContent = string(msg)

		// Broadcast the updated content to all clients.
		broadcast <- documentContent
	}
}

func handleBroadcast() {
	for {
		// Get the next message from the broadcast channel.
		msg := <-broadcast
		// Send it to every client.
		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func main() {
	// Serve static files (HTML, CSS, JavaScript).
	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/", fs)

	// WebSocket endpoint.
	http.HandleFunc("/ws", handleConnections)

	// Start the broadcaster.
	go handleBroadcast()

	// Start the server.
	fmt.Println("Server started on :8080")
	err := http.ListenAndServe("0.0.0.0:8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
