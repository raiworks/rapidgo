package jobs

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/RAiWorks/RapidGo/core/queue"
)

// RegisterJobs registers all application job handlers.
func RegisterJobs() {
	queue.RegisterHandler("example", HandleExampleJob)
}

// ExamplePayload is the payload for the example job.
type ExamplePayload struct {
	Message string `json:"message"`
}

// HandleExampleJob processes an example job.
func HandleExampleJob(_ context.Context, payload json.RawMessage) error {
	var p ExamplePayload
	if err := json.Unmarshal(payload, &p); err != nil {
		return err
	}
	slog.Info("example job processed", "message", p.Message)
	return nil
}
