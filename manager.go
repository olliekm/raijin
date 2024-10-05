package main

import (
	"log"
	"net/http"
	"raijin/auth"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	websocketUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			// Allow connections from your frontend origin
			return true
		},
	}
)

type Manager struct {
	clients      ClientList
	sync.RWMutex // stops finish line?
}

func NewManager() *Manager {
	return &Manager{
		clients: make(ClientList),
	}
}

func (m *Manager) serveWS(w http.ResponseWriter, r *http.Request) {
	log.Println("New connection")

	// Check authentication header (jwt tokens)
	token := r.Header.Get("Authorization")
	if !isAuthenticated(token) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	username, err := auth.ValidateJWT(token)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	// upgrade regular http connection to websocket
	conn, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := NewClient(conn, m, username) // Pass the username to the client
	m.addClient(client)

	// Start two go routines
	go client.readMessages()
	go client.writeMessages()

}

func (m *Manager) addClient(client *Client) {
	// Prevent multiple clients from being added at once
	m.Lock()
	// Defer will unlock the client once done adding
	defer m.Unlock()

	m.clients[client] = true
}

func (m *Manager) removeClient(client *Client) {
	// Prevent multiple clients from being added at once
	m.Lock()
	// Defer will unlock the client once done adding
	defer m.Unlock()

	if _, ok := m.clients[client]; ok {
		client.connection.Close()
		delete(m.clients, client)
	}
}

func isAuthenticated(token string) bool {
	// Add your authentication logic here
	// For example, check if the token is valid
	// This could involve checking against a database or using a JWT library
	return token == "valid_token" // Placeholder for actual token validation logic
}
