package providers

import (
	"os"
	"path/filepath"

	"github.com/RAiWorks/RGo/core/container"
	"github.com/RAiWorks/RGo/core/health"
	"github.com/RAiWorks/RGo/core/router"
	"github.com/RAiWorks/RGo/routes"
	"gorm.io/gorm"
)

// RouterProvider creates the router and registers route definitions.
type RouterProvider struct{}

// Register creates a new Router and registers it as "router" in the container.
func (p *RouterProvider) Register(c *container.Container) {
	c.Instance("router", router.New())
}

// Boot sets up the template engine, static serving, and loads route definitions.
func (p *RouterProvider) Boot(c *container.Container) {
	r := container.MustMake[*router.Router](c, "router")

	// Template engine setup — only if views directory exists
	r.SetFuncMap(router.DefaultFuncMap())
	viewsDir := filepath.Join("resources", "views")
	if info, err := os.Stat(viewsDir); err == nil && info.IsDir() {
		r.LoadTemplates(viewsDir)
	}

	// Static file serving
	if info, err := os.Stat("resources/static"); err == nil && info.IsDir() {
		r.Static("/static", "./resources/static")
	}
	if info, err := os.Stat("storage/uploads"); err == nil && info.IsDir() {
		r.Static("/uploads", "./storage/uploads")
	}

	// Route definitions
	routes.RegisterWeb(r)
	routes.RegisterAPI(r)

	// Health check endpoints (only when database is registered)
	if c.Has("db") {
		health.Routes(r, func() *gorm.DB {
			return container.MustMake[*gorm.DB](c, "db")
		})
	}
}
