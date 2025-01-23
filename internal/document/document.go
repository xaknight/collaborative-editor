package document

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Document represents a single collaborative document.
type Document struct {
	EntryCode   string
	Content     string                     // In a real implementation, this would be an OT/CRDT data structure
	Clients     map[*websocket.Conn]string // Client connection -> User's name
	mu          sync.Mutex                 // Mutex to protect concurrent access to the document
	broadcast   chan string                // Channel to broadcast updates to clients
	lastUpdated time.Time
}

// DocumentManager manages all active documents.
type DocumentManager struct {
	documents map[string]*Document
	mu        sync.RWMutex // Read/Write Mutex for concurrent access to the documents map
}

// NewDocumentManager creates a new DocumentManager.
func NewDocumentManager() *DocumentManager {
	return &DocumentManager{
		documents: make(map[string]*Document),
	}
}

// CreateDocument creates a new document with a unique entry code.
func (dm *DocumentManager) CreateDocument() *Document {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	entryCode := generateEntryCode()
	doc := &Document{
		EntryCode:   entryCode,
		Content:     "",
		Clients:     make(map[*websocket.Conn]string),
		broadcast:   make(chan string),
		lastUpdated: time.Now(),
	}
	dm.documents[entryCode] = doc
	go doc.startBroadcasting()
	return doc
}

// GetDocument retrieves a document by its entry code.
func (dm *DocumentManager) GetDocument(entryCode string) (*Document, error) {
	dm.mu.RLock()
	defer dm.mu.RUnlock()

	doc, ok := dm.documents[entryCode]
	if !ok {
		return nil, fmt.Errorf("document not found")
	}
	return doc, nil
}

// AddClient adds a client to a document.
func (doc *Document) AddClient(conn *websocket.Conn, userName string) {
	doc.mu.Lock()
	defer doc.mu.Unlock()

	doc.Clients[conn] = userName
	log.Printf("Client %s added to document %s", userName, doc.EntryCode)

	// Send the current document content to the newly connected client.
	if err := conn.WriteMessage(websocket.TextMessage, []byte(doc.Content)); err != nil {
		log.Printf("Error sending initial content to client: %v", err)
	}
}

// RemoveClient removes a client from a document.
func (doc *Document) RemoveClient(conn *websocket.Conn) {
	doc.mu.Lock()
	defer doc.mu.Unlock()

	userName, ok := doc.Clients[conn]
	if !ok {
		return // Client not found
	}

	delete(doc.Clients, conn)
	log.Printf("Client %s removed from document %s", userName, doc.EntryCode)
}

// UpdateContent updates the document's content and broadcasts the change.
func (doc *Document) UpdateContent(newContent string) {
	doc.mu.Lock()
	defer doc.mu.Unlock()

	doc.Content = newContent
	doc.lastUpdated = time.Now()
	doc.broadcast <- newContent // Send the updated content to the broadcast channel
}

// startBroadcasting handles broadcasting messages to all clients in a document.
func (doc *Document) startBroadcasting() {
	for msg := range doc.broadcast {
		doc.mu.Lock()
		for clientConn := range doc.Clients {
			err := clientConn.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				log.Printf("Error broadcasting to client: %v", err)
				clientConn.Close()
				doc.RemoveClient(clientConn)
			}
		}
		doc.mu.Unlock()
	}
}

// generateEntryCode generates a simple, random entry code (for demonstration purposes).
func generateEntryCode() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const codeLength = 6

	var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, codeLength)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func (dm *DocumentManager) DeleteDocument(entryCode string) {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	if doc, ok := dm.documents[entryCode]; ok {
		// Close all client connections for this document
		for clientConn := range doc.Clients {
			clientConn.Close()
		}

		// Stop the broadcasting goroutine
		close(doc.broadcast)

		// Remove the document from the map
		delete(dm.documents, entryCode)

		log.Printf("Document %s deleted", entryCode)
	}
}

func (dm *DocumentManager) CleanupInactiveDocuments() {
	for {
		time.Sleep(30 * time.Minute) // Check every 30 minutes

		dm.mu.Lock()
		for entryCode, doc := range dm.documents {
			if time.Since(doc.lastUpdated) > 24*time.Hour {
				dm.DeleteDocument(entryCode)
			}
		}
		dm.mu.Unlock()
	}
}
