package websocket

import (
	"context"
	"errors"
	"sync"
	"time"

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

// HubConfig controls Hub behavior.
type HubConfig struct {
	PingInterval time.Duration // How often to send Ping frames. 0 disables heartbeat. Default: 30s.
	PongTimeout  time.Duration // How long to wait for Pong after Ping. Default: 10s.
}

// Hub manages WebSocket clients and rooms.
type Hub struct {
	clients   map[string]*Client
	rooms     map[string]map[string]*Client
	mu        sync.RWMutex
	config    HubConfig
	onJoinFn  func(client *Client, room string)
	onLeaveFn func(client *Client, room string)
}

// NewHub creates an empty Hub with default config (30s ping, 10s pong timeout).
func NewHub() *Hub {
	return NewHubWithConfig(HubConfig{
		PingInterval: 30 * time.Second,
		PongTimeout:  10 * time.Second,
	})
}

// NewHubWithConfig creates an empty Hub with the given configuration.
func NewHubWithConfig(cfg HubConfig) *Hub {
	return &Hub{
		clients: make(map[string]*Client),
		rooms:   make(map[string]map[string]*Client),
		config:  cfg,
	}
}

// OnJoin registers a callback that fires when a client joins a room.
func (h *Hub) OnJoin(fn func(client *Client, room string)) {
	h.onJoinFn = fn
}

// OnLeave registers a callback that fires when a client leaves a room
// (including when removed on disconnect).
func (h *Hub) OnLeave(fn func(client *Client, room string)) {
	h.onLeaveFn = fn
}

// Handler returns a gin.HandlerFunc that upgrades the HTTP connection to
// WebSocket, creates a Client, registers it with the Hub, and calls
// onConnect. The client is automatically removed when onConnect returns.
// If PingInterval > 0, a heartbeat goroutine detects dead connections.
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

		// Start heartbeat if configured
		if h.config.PingInterval > 0 {
			ctx, cancel := context.WithCancel(c.Request.Context())
			defer cancel()
			go h.heartbeat(ctx, conn)
		}

		onConnect(client)

		h.Remove(client)
		conn.Close(websocket.StatusNormalClosure, "")
	}
}

// heartbeat sends Ping frames at the configured interval.
// If a Pong is not received within PongTimeout, the connection's context
// is cancelled, which causes the read loop in onConnect to fail and
// triggers cleanup via Remove.
func (h *Hub) heartbeat(ctx context.Context, conn *websocket.Conn) {
	ticker := time.NewTicker(h.config.PingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			pingCtx, cancel := context.WithTimeout(ctx, h.config.PongTimeout)
			err := conn.Ping(pingCtx)
			cancel()
			if err != nil {
				conn.Close(websocket.StatusGoingAway, "ping timeout")
				return
			}
		}
	}
}

// Join adds the client to the named room. The room is created if it does
// not exist. Joining the same room twice is a no-op.
func (h *Hub) Join(client *Client, room string) {
	h.mu.Lock()
	members, ok := h.rooms[room]
	if !ok {
		members = make(map[string]*Client)
		h.rooms[room] = members
	}
	_, alreadyIn := members[client.ID]
	members[client.ID] = client
	h.mu.Unlock()

	if !alreadyIn && h.onJoinFn != nil {
		h.onJoinFn(client, room)
	}
}

// Leave removes the client from the named room. The room is deleted if
// it becomes empty.
func (h *Hub) Leave(client *Client, room string) {
	h.mu.Lock()
	members, ok := h.rooms[room]
	if !ok {
		h.mu.Unlock()
		return
	}
	_, wasIn := members[client.ID]
	delete(members, client.ID)
	if len(members) == 0 {
		delete(h.rooms, room)
	}
	h.mu.Unlock()

	if wasIn && h.onLeaveFn != nil {
		h.onLeaveFn(client, room)
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
// Fires onLeaveFn for each room the client was in.
func (h *Hub) Remove(client *Client) {
	h.mu.Lock()
	delete(h.clients, client.ID)

	var leftRooms []string
	for room, members := range h.rooms {
		if _, ok := members[client.ID]; ok {
			leftRooms = append(leftRooms, room)
			delete(members, client.ID)
			if len(members) == 0 {
				delete(h.rooms, room)
			}
		}
	}
	h.mu.Unlock()

	if h.onLeaveFn != nil {
		for _, room := range leftRooms {
			h.onLeaveFn(client, room)
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
