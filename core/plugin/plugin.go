package plugin

import (
	"fmt"

	"github.com/RAiWorks/RapidGo/core/container"
	"github.com/RAiWorks/RapidGo/core/events"
	"github.com/RAiWorks/RapidGo/core/router"
	"github.com/spf13/cobra"
)

// Plugin is a self-contained module that integrates with the RapidGo framework.
// It embeds container.Provider for service registration (Register + Boot).
type Plugin interface {
	container.Provider

	// Name returns a unique identifier for the plugin (e.g., "notifications").
	Name() string
}

// RouteRegistrar is implemented by plugins that register HTTP routes.
type RouteRegistrar interface {
	RegisterRoutes(r *router.Router)
}

// CommandRegistrar is implemented by plugins that add CLI commands.
type CommandRegistrar interface {
	Commands() []*cobra.Command
}

// EventRegistrar is implemented by plugins that listen to framework events.
type EventRegistrar interface {
	RegisterEvents(d *events.Dispatcher)
}

// PluginManager coordinates plugin registration and subsystem wiring.
type PluginManager struct {
	plugins []Plugin
	names   map[string]bool
}

// NewManager creates a PluginManager.
func NewManager() *PluginManager {
	return &PluginManager{
		names: make(map[string]bool),
	}
}

// Add registers a plugin. Returns an error if a plugin with the same Name() is already registered.
func (m *PluginManager) Add(p Plugin) error {
	name := p.Name()
	if m.names[name] {
		return fmt.Errorf("plugin %q is already registered", name)
	}
	m.names[name] = true
	m.plugins = append(m.plugins, p)
	return nil
}

// Plugins returns all registered plugins.
func (m *PluginManager) Plugins() []Plugin {
	return m.plugins
}

// RegisterAll calls Register(c) on each plugin in registration order.
func (m *PluginManager) RegisterAll(c *container.Container) {
	for _, p := range m.plugins {
		p.Register(c)
	}
}

// BootAll calls Boot(c) on each plugin in registration order.
func (m *PluginManager) BootAll(c *container.Container) {
	for _, p := range m.plugins {
		p.Boot(c)
	}
}

// RegisterRoutes calls RegisterRoutes(r) on each plugin that implements RouteRegistrar.
func (m *PluginManager) RegisterRoutes(r *router.Router) {
	for _, p := range m.plugins {
		if rr, ok := p.(RouteRegistrar); ok {
			rr.RegisterRoutes(r)
		}
	}
}

// RegisterCommands calls Commands() on each plugin that implements CommandRegistrar
// and adds the returned commands to the root Cobra command.
func (m *PluginManager) RegisterCommands(root *cobra.Command) {
	for _, p := range m.plugins {
		if cr, ok := p.(CommandRegistrar); ok {
			for _, cmd := range cr.Commands() {
				root.AddCommand(cmd)
			}
		}
	}
}

// RegisterEvents calls RegisterEvents(d) on each plugin that implements EventRegistrar.
func (m *PluginManager) RegisterEvents(d *events.Dispatcher) {
	for _, p := range m.plugins {
		if er, ok := p.(EventRegistrar); ok {
			er.RegisterEvents(d)
		}
	}
}
