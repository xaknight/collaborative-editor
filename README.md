# Collaborative Code Editor

This project is a real-time collaborative code editor built using Go (Golang) on the backend and HTML, CSS, and JavaScript on the frontend. It allows multiple users to simultaneously edit code in a shared document, seeing each other's changes in real-time.

**Project Structure:**
collaborative-editor/

├── main.go // Main application

├── internal/ // Internal packages (handlers, document, user)

├── pkg/ // Public packages (api)

├── web/ // Frontend (static files, templates)

├── go.mod // Go modules

├── go.sum // Go modules checksums

└── README.md // Project description


**Features (Current & Planned):**

*   **Real-time Collaboration:** Multiple users can edit the same document simultaneously.
*   **Unique Entry Codes:** Each document has a unique, randomly generated entry code for sharing.
*   **Basic Code Editing:** Currently uses a simple text area for code input. (Planned: Integrate Ace, Monaco, or CodeMirror for a richer code editing experience).
*   **Last-Write-Wins (Temporary):** Currently, the last edit to the document overwrites previous changes. (Planned: Implement Operational Transformation (OT) or Conflict-free Replicated Data Types (CRDTs) for robust conflict resolution).
*   **Cursor Tracking (Planned):** Display the cursors of other users in real-time, along with their names.
*   **Split View (Planned):** A resizable split view with the code editor on one side and a rendered output pane on the other.
*   **Open in New Page (Planned):** Allow users to open the rendered output in a new browser tab/window.
*   **User Authentication (Planned):** Implement user accounts and authentication for secure access and document management.
*   **Persistence (Planned):** Store documents persistently using a database.

**Technology Stack:**

*   **Backend:** Go (Golang)
*   **Frontend:** HTML, CSS, JavaScript
*   **WebSocket Library:** `github.com/gorilla/websocket`
*   **Code Editor (Planned):** Ace Editor, Monaco Editor, or CodeMirror
*   **Database (Planned):** PostgreSQL, MongoDB, or similar

**Getting Started:**

1. **Prerequisites:**
    *   Go (Golang) installed on your system.
    *   A code editor or IDE.

2. **Clone the Repository:**

    ```bash
    git clone https://github.com/xaknight/collaborative-editor.git
    cd collaborative-editor
    ```

3. **Install Dependencies:**

    ```bash
    go mod tidy
    ```

4. **Run the Application:**

    ```bash
    go run main.go
    ```

5. **Access the Application:**
    *   Open your web browser and go to `http://localhost:8080`.

**How to Use:**

1. **Main Page:**
    *   Enter your name (or a random name will be generated).
    *   You can either:
        *   Create a new document by clicking "Generate Code" and then "Join/Create."
        *   Join an existing document by entering its entry code and clicking "Join/Create."
2. **Document Page:**
    *   You will be redirected to a URL like `/doc/<entry-code>`.
    *   Start typing in the text area. Changes will be reflected in real-time for other users in the same document (currently with last-write-wins).

**Current Limitations:**

*   **No OT/CRDT:** The current implementation uses a simple last-write-wins approach for conflict resolution. This can lead to data loss if multiple users edit the same part of the document simultaneously. Implementing OT/CRDTs is a high priority.
*   **Basic Text Area:** The code editor is currently a plain text area. Integrating a more advanced code editor will significantly improve the user experience.
*   **No Persistence:** Documents are only stored in memory. If the server restarts, all data is lost.
*   **No User Authentication:** There is no user management or authentication, so anyone with the entry code can access a document.

**Future Development:**

*   **Implement OT/CRDTs:** This is the most important next step to ensure robust real-time collaboration without data loss.
*   **Integrate a Code Editor:** Add a feature-rich code editor like Ace, Monaco, or CodeMirror.
*   **Add Cursor Tracking:** Display remote cursors with user names.
*   **Implement Split View:** Create a resizable split view for code editing and output rendering.
*   **Add Persistence:** Integrate a database to store documents persistently.
*   **Implement User Authentication:** Add user accounts and authentication for secure access.
*   **Improve Scalability:** Optimize the backend to handle a large number of concurrent users and documents.


**License:**

This project is licensed under the MIT License.
