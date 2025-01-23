package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"xaknight/DocCollaboration/internal/document"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow connections from any origin (for development only!)
	},
}

var docManager = document.NewDocumentManager()

// HandleJoinOrCreate handles requests to create a new document or join an existing one.
func HandleJoinOrCreate(w http.ResponseWriter, r *http.Request) {
	type JoinRequest struct {
		Name      string `json:"name"`
		EntryCode string `json:"entryCode"`
	}

	type JoinResponse struct {
		EntryCode string `json:"entryCode"`
	}

	var req JoinRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var doc *document.Document
	var err error

	if req.EntryCode == "" {
		// Create a new document
		doc = docManager.CreateDocument()
	} else {
		// Try to get an existing document
		doc, err = docManager.GetDocument(req.EntryCode)
		if err != nil {
			// If the entry code is not empty and the document doesn't exist, return an error.
			http.Error(w, "Document not found", http.StatusNotFound)
			return
		}
	}

	// Respond with the entry code
	resp := JoinResponse{EntryCode: doc.EntryCode}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// HandleConnections handles WebSocket connections for a specific document.
func HandleConnections(w http.ResponseWriter, r *http.Request) {
	entryCode := r.URL.Query().Get("doc")
	userName := r.URL.Query().Get("name")

	if entryCode == "" {
		http.Error(w, "Missing document entry code", http.StatusBadRequest)
		return
	}
	if userName == "" {
		// Generate a unique user name if not provided
		userName = "user-" + uuid.New().String()[:8]
	}

	doc, err := docManager.GetDocument(entryCode)
	if err != nil {
		http.Error(w, "Document not found", http.StatusNotFound)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading to WebSocket:", err)
		return
	}

	// Add the client to the document.
	doc.AddClient(conn, userName)

	// Handle client disconnection
	defer func() {
		doc.RemoveClient(conn)
		conn.Close()
	}()

	for {
		// Read messages from the client.
		_, msg, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break // Exit the loop if the client disconnects or an error occurs
		}

		// Update the document content (basic last-write-wins, no OT/CRDT yet).
		doc.UpdateContent(string(msg))
	}
}

func StartCleanupTask() {
	go docManager.CleanupInactiveDocuments()
}
