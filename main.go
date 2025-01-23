package main

import (
	"fmt"
	"log"
	"net/http"
	"xaknight/DocCollaboration/internal/handlers"
)

func main() {
	// Serve static files (HTML, CSS, JavaScript).
	fs := http.FileServer(http.Dir("./web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Serve HTML templates
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.ServeFile(w, r, "./web/templates/index.html")
		} else {
			http.NotFound(w, r)
		}
	})

	http.HandleFunc("/doc/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/templates/document.html")
	})

	// API endpoint to create or join a document
	http.HandleFunc("/api/join", handlers.HandleJoinOrCreate)

	// WebSocket endpoint.
	http.HandleFunc("/ws", handlers.HandleConnections)

	// Start the cleanup task
	handlers.StartCleanupTask()

	// Start the server.
	fmt.Println("Server started on :8080")
	err := http.ListenAndServe("0.0.0.0:8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
