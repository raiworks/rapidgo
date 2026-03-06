package websocket

import (
	"context"

	"github.com/coder/websocket"
	"github.com/gin-gonic/gin"
)

// Handler is a callback for handling a WebSocket connection.
type Handler func(conn *websocket.Conn, ctx context.Context)

// Options configures the WebSocket upgrade.
type Options struct {
	OriginPatterns     []string // Allowed origin patterns for CORS
	InsecureSkipVerify bool     // Skip origin verification (dev only)
}

// Upgrader returns a gin.HandlerFunc that upgrades the HTTP connection
// to WebSocket and delegates to the provided handler. If opts is nil,
// all origins are accepted (InsecureSkipVerify = true).
func Upgrader(handler Handler, opts *Options) gin.HandlerFunc {
	return func(c *gin.Context) {
		acceptOpts := &websocket.AcceptOptions{}
		if opts != nil {
			acceptOpts.OriginPatterns = opts.OriginPatterns
			acceptOpts.InsecureSkipVerify = opts.InsecureSkipVerify
		} else {
			acceptOpts.InsecureSkipVerify = true
		}

		conn, err := websocket.Accept(c.Writer, c.Request, acceptOpts)
		if err != nil {
			return
		}
		defer conn.CloseNow()

		handler(conn, c.Request.Context())

		conn.Close(websocket.StatusNormalClosure, "")
	}
}

// Echo is a built-in handler that reads messages and echoes them back.
func Echo(conn *websocket.Conn, ctx context.Context) {
	for {
		msgType, msg, err := conn.Read(ctx)
		if err != nil {
			return
		}
		if err := conn.Write(ctx, msgType, msg); err != nil {
			return
		}
	}
}
