package websocket

import (
	"context"
	"errors"
	"sync"

	"github.com/coder/websocket"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Client represents a single WebSocket connection registered with a Hub.
type Client struct {
	ID   string
	Conn *websocket.Conn
	hub  *Hub
}

// Hub manages WebSocket clients and rooms.
type Hub struct {
	clients map[string]*Client
	rooms   map[string]map[string]*Client
	mu      sync.RWMutex
}

// NewHub creates an empty Hub with initialized maps.
func NewHub() *Hub {
	return &Hub{
		clients: make(map[string]*Client),
		rooms:   make(map[string]map[string]*Client),
	}
}

// Handler returns a gin.HandlerFunc that upgrades the HTTP connection to
// WebSocket, creates a Client, registers it with the Hub, and calls
// onConnect. The client is automatically removed when onConnect returns.
func (h *Hub) Handler(onConnect func(*Client)) gin.HandlerFunc {
	return func(c *gin.Context) {
		conn, err := websocket.Accept(c.Writer, c.Request, &websocket.AcceptOptions{
			InsecureSkipVerify: true,
		})
		if err != nil {
			return
		}
		defer conn.CloseNow()

		client := &Client{
			ID:   uuid.NewString(),
			Conn: conn,
			hub:  h,
		}

		h.mu.Lock()
		h.clients[client.ID] = client
		h.mu.Unlock()

		onConnect(client)

		h.Remove(client)
		conn.Close(websocket.StatusNormalClosure, "")
	}
}

// Join adds the client to the named room. The room is created if it does
// not exist. Joining the same room twice is a no-op.
func (h *Hub) Join(client *Client, room string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	members, ok := h.rooms[room]
	if !ok {
		members = make(map[string]*Client)
		h.rooms[room] = members
	}
	members[client.ID] = client
}

// Leave removes the client from the named room. The room is deleted if
// it becomes empty.
func (h *Hub) Leave(client *Client, room string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	members, ok := h.rooms[room]
	if !ok {
		return
	}
	delete(members, client.ID)
	if len(members) == 0 {
		delete(h.rooms, room)
	}
}

// Broadcast sends a message to all clients in the named room.
func (h *Hub) Broadcast(room string, msgType websocket.MessageType, data []byte) {
	h.mu.RLock()
	members := h.rooms[room]
	targets := make([]*Client, 0, len(members))
	for _, c := range members {
		targets = append(targets, c)
	}
	h.mu.RUnlock()

	for _, c := range targets {
		c.Conn.Write(context.Background(), msgType, data)
	}
}

// BroadcastOthers sends a message to all clients in the named room except
// the sender.
func (h *Hub) BroadcastOthers(room string, sender *Client, msgType websocket.MessageType, data []byte) {
	h.mu.RLock()
	members := h.rooms[room]
	targets := make([]*Client, 0, len(members))
	for _, c := range members {
		if c.ID != sender.ID {
			targets = append(targets, c)
		}
	}
	h.mu.RUnlock()

	for _, c := range targets {
		c.Conn.Write(context.Background(), msgType, data)
	}
}

// Send delivers a message to a single client by ID. Returns an error if
// the client is not found.
func (h *Hub) Send(clientID string, msgType websocket.MessageType, data []byte) error {
	h.mu.RLock()
	client, ok := h.clients[clientID]
	h.mu.RUnlock()

	if !ok {
		return errors.New("client not found")
	}
	return client.Conn.Write(context.Background(), msgType, data)
}

// Remove removes the client from all rooms and from the Hub.
func (h *Hub) Remove(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	delete(h.clients, client.ID)

	for room, members := range h.rooms {
		delete(members, client.ID)
		if len(members) == 0 {
			delete(h.rooms, room)
		}
	}
}

// Clients returns a snapshot of all clients in the named room.
// Returns nil if the room does not exist.
func (h *Hub) Clients(room string) []*Client {
	h.mu.RLock()
	defer h.mu.RUnlock()

	members, ok := h.rooms[room]
	if !ok {
		return nil
	}
	result := make([]*Client, 0, len(members))
	for _, c := range members {
		result = append(result, c)
	}
	return result
}

// Rooms returns a snapshot of all active room names.
func (h *Hub) Rooms() []string {
	h.mu.RLock()
	defer h.mu.RUnlock()

	result := make([]string, 0, len(h.rooms))
	for name := range h.rooms {
		result = append(result, name)
	}
	return result
}
