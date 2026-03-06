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

func init() {
	gin.SetMode(gin.TestMode)
}

// newWSServer creates a test server with the given handler and options.
func newWSServer(handler Handler, opts *Options) *httptest.Server {
	e := gin.New()
	e.GET("/ws", Upgrader(handler, opts))
	return httptest.NewServer(e)
}

// wsURL converts http://... to ws://...
func wsURL(s *httptest.Server) string {
	return "ws" + strings.TrimPrefix(s.URL, "http") + "/ws"
}

// TC-01: Upgrader upgrades HTTP to WebSocket
func TestUpgrader_Connects(t *testing.T) {
	s := newWSServer(Echo, nil)
	defer s.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, _, err := websocket.Dial(ctx, wsURL(s), nil)
	if err != nil {
		t.Fatalf("Dial failed: %v", err)
	}
	conn.Close(websocket.StatusNormalClosure, "")
}

// TC-02: Echo handler echoes text message
func TestEcho_TextMessage(t *testing.T) {
	s := newWSServer(Echo, nil)
	defer s.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, _, err := websocket.Dial(ctx, wsURL(s), nil)
	if err != nil {
		t.Fatalf("Dial failed: %v", err)
	}
	defer conn.CloseNow()

	msg := "hello websocket"
	if err := conn.Write(ctx, websocket.MessageText, []byte(msg)); err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	typ, data, err := conn.Read(ctx)
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}
	if typ != websocket.MessageText {
		t.Fatalf("expected MessageText, got %v", typ)
	}
	if string(data) != msg {
		t.Fatalf("expected %q, got %q", msg, string(data))
	}
}

// TC-03: Echo handler echoes binary message
func TestEcho_BinaryMessage(t *testing.T) {
	s := newWSServer(Echo, nil)
	defer s.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, _, err := websocket.Dial(ctx, wsURL(s), nil)
	if err != nil {
		t.Fatalf("Dial failed: %v", err)
	}
	defer conn.CloseNow()

	payload := []byte{0x01, 0x02, 0x03, 0xFF}
	if err := conn.Write(ctx, websocket.MessageBinary, payload); err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	typ, data, err := conn.Read(ctx)
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}
	if typ != websocket.MessageBinary {
		t.Fatalf("expected MessageBinary, got %v", typ)
	}
	if len(data) != len(payload) {
		t.Fatalf("expected %d bytes, got %d", len(payload), len(data))
	}
	for i, b := range data {
		if b != payload[i] {
			t.Fatalf("byte %d: expected %x, got %x", i, payload[i], b)
		}
	}
}

// TC-04: Connection closes cleanly on handler return
func TestUpgrader_CleanClose(t *testing.T) {
	// Handler that returns immediately — server should close the connection.
	s := newWSServer(func(conn *websocket.Conn, ctx context.Context) {
		// do nothing, return immediately
	}, nil)
	defer s.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, _, err := websocket.Dial(ctx, wsURL(s), nil)
	if err != nil {
		t.Fatalf("Dial failed: %v", err)
	}
	defer conn.CloseNow()

	// Read should return an error (server closed)
	_, _, err = conn.Read(ctx)
	if err == nil {
		t.Fatal("expected error after server close")
	}
}

// TC-05: Upgrader with nil options uses defaults
func TestUpgrader_NilOptions(t *testing.T) {
	s := newWSServer(Echo, nil)
	defer s.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, _, err := websocket.Dial(ctx, wsURL(s), nil)
	if err != nil {
		t.Fatalf("Dial with nil options failed: %v", err)
	}
	conn.Close(websocket.StatusNormalClosure, "")
}

// TC-06: Custom handler receives connection
func TestUpgrader_CustomHandler(t *testing.T) {
	var mu sync.Mutex
	called := false

	handler := func(conn *websocket.Conn, ctx context.Context) {
		mu.Lock()
		called = true
		mu.Unlock()
	}

	s := newWSServer(handler, nil)
	defer s.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, _, err := websocket.Dial(ctx, wsURL(s), nil)
	if err != nil {
		t.Fatalf("Dial failed: %v", err)
	}
	defer conn.CloseNow()

	// Give server time to invoke the handler
	time.Sleep(100 * time.Millisecond)

	mu.Lock()
	if !called {
		t.Fatal("custom handler was not invoked")
	}
	mu.Unlock()
}
