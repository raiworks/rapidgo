package plugin

import (
	"testing"

	"github.com/RAiWorks/RapidGo/core/container"
	"github.com/RAiWorks/RapidGo/core/events"
	"github.com/RAiWorks/RapidGo/core/router"
	"github.com/spf13/cobra"
)

// --- Mock Plugins ---

// mockPlugin implements Plugin (Provider + Name) only.
type mockPlugin struct {
	name             string
	registerCalled   bool
	bootCalled       bool
	lifecycleTracker *[]string // shared tracker for ordering tests
}

func (p *mockPlugin) Name() string { return p.name }
func (p *mockPlugin) Register(c *container.Container) {
	p.registerCalled = true
	if p.lifecycleTracker != nil {
		*p.lifecycleTracker = append(*p.lifecycleTracker, p.name+".Register")
	}
}
func (p *mockPlugin) Boot(c *container.Container) {
	p.bootCalled = true
	if p.lifecycleTracker != nil {
		*p.lifecycleTracker = append(*p.lifecycleTracker, p.name+".Boot")
	}
}

// mockRoutePlugin implements Plugin + RouteRegistrar.
type mockRoutePlugin struct {
	mockPlugin
	routesCalled bool
}

func (p *mockRoutePlugin) RegisterRoutes(r *router.Router) {
	p.routesCalled = true
}

// mockCommandPlugin implements Plugin + CommandRegistrar.
type mockCommandPlugin struct {
	mockPlugin
	commandsCalled bool
	cmds           []*cobra.Command
}

func (p *mockCommandPlugin) Commands() []*cobra.Command {
	p.commandsCalled = true
	return p.cmds
}

// mockEventPlugin implements Plugin + EventRegistrar.
type mockEventPlugin struct {
	mockPlugin
	eventsCalled bool
}

func (p *mockEventPlugin) RegisterEvents(d *events.Dispatcher) {
	p.eventsCalled = true
}

// mockFullPlugin implements Plugin + all optional interfaces.
type mockFullPlugin struct {
	mockPlugin
	routesCalled    bool
	commandsCalled  bool
	eventsCalled    bool
	cmds            []*cobra.Command
}

func (p *mockFullPlugin) RegisterRoutes(r *router.Router) { p.routesCalled = true }
func (p *mockFullPlugin) Commands() []*cobra.Command {
	p.commandsCalled = true
	return p.cmds
}
func (p *mockFullPlugin) RegisterEvents(d *events.Dispatcher) { p.eventsCalled = true }

// mockServicePlugin binds a service in Register.
type mockServicePlugin struct {
	mockPlugin
}

func (p *mockServicePlugin) Register(c *container.Container) {
	p.registerCalled = true
	c.Singleton("test.greeting", func(c *container.Container) interface{} {
		return "Hello from plugin!"
	})
}

// --- Tests ---

func TestNewManager(t *testing.T) {
	m := NewManager()
	if m == nil {
		t.Fatal("NewManager() returned nil")
	}
	if len(m.Plugins()) != 0 {
		t.Fatalf("expected 0 plugins, got %d", len(m.Plugins()))
	}
}

func TestAddPlugin(t *testing.T) {
	m := NewManager()
	err := m.Add(&mockPlugin{name: "alpha"})
	if err != nil {
		t.Fatalf("Add() returned error: %v", err)
	}
	if len(m.Plugins()) != 1 {
		t.Fatalf("expected 1 plugin, got %d", len(m.Plugins()))
	}
	if m.Plugins()[0].Name() != "alpha" {
		t.Errorf("expected name 'alpha', got %q", m.Plugins()[0].Name())
	}
}

func TestAddMultiplePlugins(t *testing.T) {
	m := NewManager()
	names := []string{"alpha", "beta", "gamma"}
	for _, n := range names {
		if err := m.Add(&mockPlugin{name: n}); err != nil {
			t.Fatalf("Add(%s) error: %v", n, err)
		}
	}
	if len(m.Plugins()) != 3 {
		t.Fatalf("expected 3 plugins, got %d", len(m.Plugins()))
	}
	for i, p := range m.Plugins() {
		if p.Name() != names[i] {
			t.Errorf("plugin %d: expected %q, got %q", i, names[i], p.Name())
		}
	}
}

func TestAddDuplicateNameReturnsError(t *testing.T) {
	m := NewManager()
	m.Add(&mockPlugin{name: "dup"})
	err := m.Add(&mockPlugin{name: "dup"})
	if err == nil {
		t.Fatal("expected error for duplicate plugin name, got nil")
	}
}

func TestAddDuplicateDoesNotModifyList(t *testing.T) {
	m := NewManager()
	m.Add(&mockPlugin{name: "dup"})
	m.Add(&mockPlugin{name: "dup"}) // second should fail
	if len(m.Plugins()) != 1 {
		t.Fatalf("expected 1 plugin after duplicate, got %d", len(m.Plugins()))
	}
}

func TestRegisterAll(t *testing.T) {
	m := NewManager()
	p1 := &mockPlugin{name: "a"}
	p2 := &mockPlugin{name: "b"}
	m.Add(p1)
	m.Add(p2)

	c := container.New()
	m.RegisterAll(c)

	if !p1.registerCalled {
		t.Error("plugin 'a' Register() not called")
	}
	if !p2.registerCalled {
		t.Error("plugin 'b' Register() not called")
	}
}

func TestBootAll(t *testing.T) {
	m := NewManager()
	p1 := &mockPlugin{name: "a"}
	p2 := &mockPlugin{name: "b"}
	m.Add(p1)
	m.Add(p2)

	c := container.New()
	m.BootAll(c)

	if !p1.bootCalled {
		t.Error("plugin 'a' Boot() not called")
	}
	if !p2.bootCalled {
		t.Error("plugin 'b' Boot() not called")
	}
}

func TestLifecycleOrder(t *testing.T) {
	tracker := make([]string, 0)
	m := NewManager()
	m.Add(&mockPlugin{name: "first", lifecycleTracker: &tracker})
	m.Add(&mockPlugin{name: "second", lifecycleTracker: &tracker})

	c := container.New()
	m.RegisterAll(c)
	m.BootAll(c)

	expected := []string{
		"first.Register", "second.Register",
		"first.Boot", "second.Boot",
	}
	if len(tracker) != len(expected) {
		t.Fatalf("expected %d lifecycle calls, got %d: %v", len(expected), len(tracker), tracker)
	}
	for i, v := range expected {
		if tracker[i] != v {
			t.Errorf("lifecycle[%d]: expected %q, got %q", i, v, tracker[i])
		}
	}
}

func TestPluginServiceAccessible(t *testing.T) {
	m := NewManager()
	m.Add(&mockServicePlugin{mockPlugin: mockPlugin{name: "svc"}})

	c := container.New()
	m.RegisterAll(c)

	greeting := container.MustMake[string](c, "test.greeting")
	if greeting != "Hello from plugin!" {
		t.Errorf("expected 'Hello from plugin!', got %q", greeting)
	}
}

func TestRegisterRoutes(t *testing.T) {
	m := NewManager()
	rp := &mockRoutePlugin{mockPlugin: mockPlugin{name: "routes"}}
	m.Add(rp)

	r := router.New()
	m.RegisterRoutes(r)

	if !rp.routesCalled {
		t.Error("RegisterRoutes() not called on RouteRegistrar plugin")
	}
}

func TestRegisterRoutesSkipsNonRegistrar(t *testing.T) {
	m := NewManager()
	plain := &mockPlugin{name: "plain"}
	m.Add(plain)

	r := router.New()
	m.RegisterRoutes(r) // should not panic
}

func TestRegisterRoutesMixed(t *testing.T) {
	m := NewManager()
	plain := &mockPlugin{name: "plain"}
	rp := &mockRoutePlugin{mockPlugin: mockPlugin{name: "routes"}}
	m.Add(plain)
	m.Add(rp)

	r := router.New()
	m.RegisterRoutes(r)

	if !rp.routesCalled {
		t.Error("RouteRegistrar plugin not called")
	}
}

func TestRegisterCommands(t *testing.T) {
	m := NewManager()
	cmd := &cobra.Command{Use: "test:cmd", Short: "test command"}
	cp := &mockCommandPlugin{
		mockPlugin: mockPlugin{name: "cmds"},
		cmds:       []*cobra.Command{cmd},
	}
	m.Add(cp)

	root := &cobra.Command{Use: "root"}
	m.RegisterCommands(root)

	if !cp.commandsCalled {
		t.Error("Commands() not called on CommandRegistrar plugin")
	}
}

func TestRegisterCommandsSkipsNonRegistrar(t *testing.T) {
	m := NewManager()
	m.Add(&mockPlugin{name: "plain"})

	root := &cobra.Command{Use: "root"}
	m.RegisterCommands(root) // should not panic
}

func TestRegisteredCommandIsUsable(t *testing.T) {
	m := NewManager()
	cmd := &cobra.Command{Use: "plugin:hello", Short: "hello from plugin"}
	m.Add(&mockCommandPlugin{
		mockPlugin: mockPlugin{name: "hello"},
		cmds:       []*cobra.Command{cmd},
	})

	root := &cobra.Command{Use: "root"}
	m.RegisterCommands(root)

	found := false
	for _, c := range root.Commands() {
		if c.Use == "plugin:hello" {
			found = true
			break
		}
	}
	if !found {
		t.Error("plugin command 'plugin:hello' not found in root commands")
	}
}

func TestRegisterEvents(t *testing.T) {
	m := NewManager()
	ep := &mockEventPlugin{mockPlugin: mockPlugin{name: "events"}}
	m.Add(ep)

	d := events.NewDispatcher()
	m.RegisterEvents(d)

	if !ep.eventsCalled {
		t.Error("RegisterEvents() not called on EventRegistrar plugin")
	}
}

func TestRegisterEventsSkipsNonRegistrar(t *testing.T) {
	m := NewManager()
	m.Add(&mockPlugin{name: "plain"})

	d := events.NewDispatcher()
	m.RegisterEvents(d) // should not panic
}

func TestFullPluginAllInterfaces(t *testing.T) {
	m := NewManager()
	cmd := &cobra.Command{Use: "full:cmd"}
	fp := &mockFullPlugin{
		mockPlugin: mockPlugin{name: "full"},
		cmds:       []*cobra.Command{cmd},
	}
	m.Add(fp)

	c := container.New()
	m.RegisterAll(c)
	m.BootAll(c)

	r := router.New()
	m.RegisterRoutes(r)

	root := &cobra.Command{Use: "root"}
	m.RegisterCommands(root)

	d := events.NewDispatcher()
	m.RegisterEvents(d)

	if !fp.registerCalled {
		t.Error("Register() not called")
	}
	if !fp.bootCalled {
		t.Error("Boot() not called")
	}
	if !fp.routesCalled {
		t.Error("RegisterRoutes() not called")
	}
	if !fp.commandsCalled {
		t.Error("Commands() not called")
	}
	if !fp.eventsCalled {
		t.Error("RegisterEvents() not called")
	}
}
