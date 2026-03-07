package websocket

import (
	"context"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/coder/websocket"
	"github.com/gin-gonic/gin"
)

// newHubServer creates a test server with Hub.Handler and the given onConnect callback.
func newHubServer(hub *Hub, onConnect func(*Client)) *httptest.Server {
	e := gin.New()
	e.GET("/ws", hub.Handler(onConnect))
	return httptest.NewServer(e)
}

// hubURL converts http://... to ws://.../ws
func hubURL(s *httptest.Server) string {
	return "ws" + strings.TrimPrefix(s.URL, "http") + "/ws"
}

// dialHub dials the hub test server and returns a connection.
func dialHub(t *testing.T, s *httptest.Server) *websocket.Conn {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	conn, _, err := websocket.Dial(ctx, hubURL(s), nil)
	if err != nil {
		t.Fatalf("Dial failed: %v", err)
	}
	return conn
}

// T01: NewHub returns non-nil Hub
func TestNewHub(t *testing.T) {
	hub := NewHub()
	if hub == nil {
		t.Fatal("NewHub returned nil")
	}
}

// T02: Handler upgrades and calls onConnect
func TestHub_Handler_Connects(t *testing.T) {
	hub := NewHub()
	var mu sync.Mutex
	called := false

	s := newHubServer(hub, func(c *Client) {
		mu.Lock()
		called = true
		mu.Unlock()
	})
	defer s.Close()

	conn := dialHub(t, s)
	defer conn.CloseNow()

	// Give time for onConnect to fire and return
	time.Sleep(100 * time.Millisecond)

	mu.Lock()
	if !called {
		t.Fatal("onConnect was not called")
	}
	mu.Unlock()
}

// T03: Client receives a non-empty ID
func TestHub_Handler_AssignsClientID(t *testing.T) {
	hub := NewHub()
	idCh := make(chan string, 1)

	s := newHubServer(hub, func(c *Client) {
		idCh <- c.ID
	})
	defer s.Close()

	conn := dialHub(t, s)
	defer conn.CloseNow()

	select {
	case id := <-idCh:
		if id == "" {
			t.Fatal("client ID is empty")
		}
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for client ID")
	}
}

// T04: Client removed from Hub when connection closes
func TestHub_Handler_RemovesOnDisconnect(t *testing.T) {
	hub := NewHub()
	done := make(chan struct{})

	s := newHubServer(hub, func(c *Client) {
		// Block until we've verified the client is registered
		<-done
	})
	defer s.Close()

	conn := dialHub(t, s)

	// Client should be registered now
	time.Sleep(100 * time.Millisecond)
	hub.mu.RLock()
	count := len(hub.clients)
	hub.mu.RUnlock()
	if count != 1 {
		t.Fatalf("expected 1 client, got %d", count)
	}

	// Let onConnect return, which triggers Remove
	close(done)
	conn.Close(websocket.StatusNormalClosure, "")

	// Wait for cleanup
	time.Sleep(200 * time.Millisecond)
	hub.mu.RLock()
	count = len(hub.clients)
	hub.mu.RUnlock()
	if count != 0 {
		t.Fatalf("expected 0 clients after disconnect, got %d", count)
	}
}

// T05: Client joins a room
func TestHub_Join(t *testing.T) {
	hub := NewHub()
	readyCh := make(chan struct{})
	done := make(chan struct{})

	s := newHubServer(hub, func(c *Client) {
		hub.Join(c, "lobby")
		close(readyCh)
		<-done
	})
	defer s.Close()

	conn := dialHub(t, s)
	defer conn.CloseNow()

	<-readyCh
	clients := hub.Clients("lobby")
	if len(clients) != 1 {
		t.Fatalf("expected 1 client in lobby, got %d", len(clients))
	}
	close(done)
}

// T06: Joining a room creates it
func TestHub_Join_CreatesRoom(t *testing.T) {
	hub := NewHub()
	readyCh := make(chan struct{})
	done := make(chan struct{})

	s := newHubServer(hub, func(c *Client) {
		hub.Join(c, "new-room")
		close(readyCh)
		<-done
	})
	defer s.Close()

	conn := dialHub(t, s)
	defer conn.CloseNow()

	<-readyCh
	rooms := hub.Rooms()
	found := false
	for _, r := range rooms {
		if r == "new-room" {
			found = true
		}
	}
	if !found {
		t.Fatal("room 'new-room' not found in Rooms()")
	}
	close(done)
}

// T07: Joining same room twice is a no-op
func TestHub_Join_Idempotent(t *testing.T) {
	hub := NewHub()
	readyCh := make(chan struct{})
	done := make(chan struct{})

	s := newHubServer(hub, func(c *Client) {
		hub.Join(c, "room")
		hub.Join(c, "room") // second join
		close(readyCh)
		<-done
	})
	defer s.Close()

	conn := dialHub(t, s)
	defer conn.CloseNow()

	<-readyCh
	clients := hub.Clients("room")
	if len(clients) != 1 {
		t.Fatalf("expected 1 client after double join, got %d", len(clients))
	}
	close(done)
}

// T08: Client leaves a room
func TestHub_Leave(t *testing.T) {
	hub := NewHub()
	readyCh := make(chan struct{})
	done := make(chan struct{})

	s := newHubServer(hub, func(c *Client) {
		hub.Join(c, "room")
		hub.Leave(c, "room")
		close(readyCh)
		<-done
	})
	defer s.Close()

	conn := dialHub(t, s)
	defer conn.CloseNow()

	<-readyCh
	clients := hub.Clients("room")
	if clients != nil {
		t.Fatalf("expected nil clients after leave, got %d", len(clients))
	}
	close(done)
}

// T09: Empty room removed after last client leaves
func TestHub_Leave_RemovesEmptyRoom(t *testing.T) {
	hub := NewHub()
	readyCh := make(chan struct{})
	done := make(chan struct{})

	s := newHubServer(hub, func(c *Client) {
		hub.Join(c, "temp-room")
		hub.Leave(c, "temp-room")
		close(readyCh)
		<-done
	})
	defer s.Close()

	conn := dialHub(t, s)
	defer conn.CloseNow()

	<-readyCh
	rooms := hub.Rooms()
	for _, r := range rooms {
		if r == "temp-room" {
			t.Fatal("temp-room should have been removed")
		}
	}
	close(done)
}

// T10: Broadcast sends to all clients in room
func TestHub_Broadcast(t *testing.T) {
	hub := NewHub()
	ready := make(chan struct{}, 2)
	done := make(chan struct{})

	s := newHubServer(hub, func(c *Client) {
		hub.Join(c, "chat")
		ready <- struct{}{}
		<-done
	})
	defer s.Close()

	conn1 := dialHub(t, s)
	defer conn1.CloseNow()
	conn2 := dialHub(t, s)
	defer conn2.CloseNow()

	// Wait for both clients to join
	<-ready
	<-ready

	hub.Broadcast("chat", websocket.MessageText, []byte("hello"))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Both clients should receive the message
	for _, conn := range []*websocket.Conn{conn1, conn2} {
		typ, data, err := conn.Read(ctx)
		if err != nil {
			t.Fatalf("Read failed: %v", err)
		}
		if typ != websocket.MessageText {
			t.Fatalf("expected MessageText, got %v", typ)
		}
		if string(data) != "hello" {
			t.Fatalf("expected 'hello', got %q", string(data))
		}
	}
	close(done)
}

// T11: BroadcastOthers skips sender
func TestHub_BroadcastOthers(t *testing.T) {
	hub := NewHub()
	ready := make(chan *Client, 2)
	done := make(chan struct{})

	s := newHubServer(hub, func(c *Client) {
		hub.Join(c, "chat")
		ready <- c
		<-done
	})
	defer s.Close()

	conn1 := dialHub(t, s)
	defer conn1.CloseNow()
	conn2 := dialHub(t, s)
	defer conn2.CloseNow()

	sender := <-ready
	<-ready // second client ready

	hub.BroadcastOthers("chat", sender, websocket.MessageText, []byte("from-sender"))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// conn2 (not the sender) should receive the message
	// We need to figure out which conn maps to which client.
	// Since sender is the first dial (conn1), conn2 should get the message.
	typ, data, err := conn2.Read(ctx)
	if err != nil {
		t.Fatalf("Read on non-sender failed: %v", err)
	}
	if typ != websocket.MessageText || string(data) != "from-sender" {
		t.Fatalf("unexpected message: type=%v data=%q", typ, string(data))
	}

	// conn1 (sender) should NOT receive — use a short timeout
	shortCtx, shortCancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer shortCancel()
	_, _, err = conn1.Read(shortCtx)
	if err == nil {
		t.Fatal("sender should not have received the message")
	}
	close(done)
}

// T12: Broadcast to room A doesn't reach room B
func TestHub_Broadcast_IgnoresOtherRooms(t *testing.T) {
	hub := NewHub()
	readyA := make(chan struct{})
	readyB := make(chan struct{})
	done := make(chan struct{})

	e := gin.New()
	e.GET("/ws/a", hub.Handler(func(c *Client) {
		hub.Join(c, "roomA")
		close(readyA)
		<-done
	}))
	e.GET("/ws/b", hub.Handler(func(c *Client) {
		hub.Join(c, "roomB")
		close(readyB)
		<-done
	}))
	s := httptest.NewServer(e)
	defer s.Close()

	urlA := "ws" + strings.TrimPrefix(s.URL, "http") + "/ws/a"
	urlB := "ws" + strings.TrimPrefix(s.URL, "http") + "/ws/b"

	ctx := context.Background()
	connA, _, err := websocket.Dial(ctx, urlA, nil)
	if err != nil {
		t.Fatalf("Dial A failed: %v", err)
	}
	defer connA.CloseNow()

	connB, _, err := websocket.Dial(ctx, urlB, nil)
	if err != nil {
		t.Fatalf("Dial B failed: %v", err)
	}
	defer connB.CloseNow()

	<-readyA
	<-readyB

	hub.Broadcast("roomA", websocket.MessageText, []byte("only-for-A"))

	readCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// connA should receive
	typ, data, err := connA.Read(readCtx)
	if err != nil {
		t.Fatalf("Read A failed: %v", err)
	}
	if typ != websocket.MessageText || string(data) != "only-for-A" {
		t.Fatalf("unexpected message on A: type=%v data=%q", typ, string(data))
	}

	// connB should NOT receive
	shortCtx, shortCancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer shortCancel()
	_, _, err = connB.Read(shortCtx)
	if err == nil {
		t.Fatal("connB in roomB should not have received roomA broadcast")
	}
	close(done)
}

// T13: Send delivers to specific client
func TestHub_Send(t *testing.T) {
	hub := NewHub()
	ready := make(chan *Client, 2)
	done := make(chan struct{})

	s := newHubServer(hub, func(c *Client) {
		ready <- c
		<-done
	})
	defer s.Close()

	conn1 := dialHub(t, s)
	defer conn1.CloseNow()
	conn2 := dialHub(t, s)
	defer conn2.CloseNow()

	target := <-ready
	<-ready // second client ready

	err := hub.Send(target.ID, websocket.MessageText, []byte("direct"))
	if err != nil {
		t.Fatalf("Send failed: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// conn1 (target) should receive
	typ, data, err := conn1.Read(ctx)
	if err != nil {
		t.Fatalf("Read on target failed: %v", err)
	}
	if typ != websocket.MessageText || string(data) != "direct" {
		t.Fatalf("unexpected message: type=%v data=%q", typ, string(data))
	}

	// conn2 should NOT receive
	shortCtx, shortCancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer shortCancel()
	_, _, err = conn2.Read(shortCtx)
	if err == nil {
		t.Fatal("non-target should not have received the message")
	}
	close(done)
}

// T14: Send returns error for unknown ID
func TestHub_Send_NotFound(t *testing.T) {
	hub := NewHub()
	err := hub.Send("nonexistent", websocket.MessageText, []byte("data"))
	if err == nil {
		t.Fatal("expected error for unknown client ID")
	}
}

// T15: Remove cleans client from all rooms
func TestHub_Remove(t *testing.T) {
	hub := NewHub()
	readyCh := make(chan *Client, 1)
	done := make(chan struct{})

	s := newHubServer(hub, func(c *Client) {
		hub.Join(c, "room1")
		hub.Join(c, "room2")
		readyCh <- c
		<-done
	})
	defer s.Close()

	conn := dialHub(t, s)
	defer conn.CloseNow()

	client := <-readyCh

	// Verify client is in both rooms
	if len(hub.Clients("room1")) != 1 {
		t.Fatal("expected client in room1")
	}
	if len(hub.Clients("room2")) != 1 {
		t.Fatal("expected client in room2")
	}

	hub.Remove(client)

	// After removal, both rooms should be empty/gone
	if hub.Clients("room1") != nil {
		t.Fatal("room1 should be empty after Remove")
	}
	if hub.Clients("room2") != nil {
		t.Fatal("room2 should be empty after Remove")
	}
	close(done)
}

// T16: Clients returns a copy, not the internal map
func TestHub_Clients_Snapshot(t *testing.T) {
	hub := NewHub()
	readyCh := make(chan struct{})
	done := make(chan struct{})

	s := newHubServer(hub, func(c *Client) {
		hub.Join(c, "snap")
		close(readyCh)
		<-done
	})
	defer s.Close()

	conn := dialHub(t, s)
	defer conn.CloseNow()

	<-readyCh
	snapshot := hub.Clients("snap")
	if len(snapshot) != 1 {
		t.Fatalf("expected 1 client, got %d", len(snapshot))
	}

	// Modify the returned slice — hub should be unaffected
	snapshot[0] = nil
	actual := hub.Clients("snap")
	if len(actual) != 1 || actual[0] == nil {
		t.Fatal("modifying snapshot should not affect hub data")
	}
	close(done)
}

// T17: Rooms returns all active room names
func TestHub_Rooms_Snapshot(t *testing.T) {
	hub := NewHub()
	readyCh := make(chan struct{})
	done := make(chan struct{})

	s := newHubServer(hub, func(c *Client) {
		hub.Join(c, "alpha")
		hub.Join(c, "beta")
		hub.Join(c, "gamma")
		close(readyCh)
		<-done
	})
	defer s.Close()

	conn := dialHub(t, s)
	defer conn.CloseNow()

	<-readyCh
	rooms := hub.Rooms()
	if len(rooms) != 3 {
		t.Fatalf("expected 3 rooms, got %d: %v", len(rooms), rooms)
	}

	expected := map[string]bool{"alpha": true, "beta": true, "gamma": true}
	for _, r := range rooms {
		if !expected[r] {
			t.Fatalf("unexpected room: %s", r)
		}
	}
	close(done)
}

// T18: Concurrent join/leave/broadcast is safe (run with -race)
func TestHub_Concurrent(t *testing.T) {
	hub := NewHub()
	ready := make(chan *Client, 4)
	done := make(chan struct{})

	s := newHubServer(hub, func(c *Client) {
		ready <- c
		<-done
	})
	defer s.Close()

	// Connect 4 clients
	conns := make([]*websocket.Conn, 4)
	for i := 0; i < 4; i++ {
		conns[i] = dialHub(t, s)
		defer conns[i].CloseNow()
	}

	clients := make([]*Client, 4)
	for i := 0; i < 4; i++ {
		clients[i] = <-ready
	}

	var wg sync.WaitGroup
	const iterations = 50

	// Goroutine: rapid join/leave
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func(c *Client) {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				hub.Join(c, "stress")
				hub.Clients("stress")
				hub.Rooms()
				hub.Leave(c, "stress")
			}
		}(clients[i])
	}

	// Goroutine: rapid broadcast
	wg.Add(1)
	go func() {
		defer wg.Done()
		for j := 0; j < iterations; j++ {
			hub.Broadcast("stress", websocket.MessageText, []byte("ping"))
		}
	}()

	wg.Wait()
	close(done)
}
