package example

import (
	"fmt"
	"log/slog"

	"github.com/RAiWorks/RapidGo/core/container"
	"github.com/RAiWorks/RapidGo/core/router"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

// ExamplePlugin demonstrates the RapidGo plugin interface.
// It implements Plugin, RouteRegistrar, and CommandRegistrar.
type ExamplePlugin struct{}

// New creates an ExamplePlugin.
func New() *ExamplePlugin { return &ExamplePlugin{} }

// --- Plugin interface ---

// Name returns the plugin's unique identifier.
func (p *ExamplePlugin) Name() string { return "example" }

// Register binds plugin services into the container.
func (p *ExamplePlugin) Register(c *container.Container) {
	c.Singleton("example.greeting", func(c *container.Container) interface{} {
		return "Hello from the Example Plugin!"
	})
}

// Boot runs after all providers and plugins have been registered.
func (p *ExamplePlugin) Boot(c *container.Container) {
	slog.Info("example plugin booted")
}

// --- RouteRegistrar ---

// RegisterRoutes adds plugin routes to the router.
func (p *ExamplePlugin) RegisterRoutes(r *router.Router) {
	g := r.Group("/example")
	g.Get("/hello", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Hello from ExamplePlugin"})
	})
}

// --- CommandRegistrar ---

// Commands returns CLI commands provided by the plugin.
func (p *ExamplePlugin) Commands() []*cobra.Command {
	return []*cobra.Command{
		{
			Use:   "example:greet",
			Short: "Print a greeting from the example plugin",
			Run: func(cmd *cobra.Command, args []string) {
				fmt.Println("Hello from ExamplePlugin CLI!")
			},
		},
	}
}
