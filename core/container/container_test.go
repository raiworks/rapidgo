package container

import (
	"sync"
	"sync/atomic"
	"testing"
)

// TC-01: Bind registers transient factory — new instance each time
func TestBind_TransientFactory(t *testing.T) {
	c := New()
	var counter int
	c.Bind("counter", func(_ *Container) interface{} {
		counter++
		return counter
	})

	first := c.Make("counter").(int)
	second := c.Make("counter").(int)

	if first == second {
		t.Errorf("Bind should create new instance each time, got same value: %d", first)
	}
	if first != 1 || second != 2 {
		t.Errorf("expected 1 and 2, got %d and %d", first, second)
	}
}

// TC-02: Singleton creates instance only once
func TestSingleton_CreatesOnce(t *testing.T) {
	c := New()
	var callCount int
	c.Singleton("db", func(_ *Container) interface{} {
		callCount++
		return "db-instance"
	})

	r1 := c.Make("db")
	r2 := c.Make("db")
	r3 := c.Make("db")

	if callCount != 1 {
		t.Errorf("factory should be called once, called %d times", callCount)
	}
	if r1 != r2 || r2 != r3 {
		t.Error("all calls should return the same instance")
	}
}

// TC-03: Instance registers pre-created object
func TestInstance_PreCreated(t *testing.T) {
	c := New()
	type Config struct{ Name string }
	cfg := &Config{Name: "test"}
	c.Instance("config", cfg)

	result := c.Make("config").(*Config)
	if result != cfg {
		t.Error("should return the exact same pointer")
	}
	if result.Name != "test" {
		t.Errorf("expected 'test', got '%s'", result.Name)
	}
}

// TC-04: Make panics on unregistered service
func TestMake_PanicsOnMissing(t *testing.T) {
	c := New()

	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic on unregistered service")
		}
		msg, ok := r.(string)
		if !ok || msg != "service not found: nonexistent" {
			t.Errorf("unexpected panic message: %v", r)
		}
	}()

	c.Make("nonexistent")
}

// TC-05: MustMake resolves with correct type
func TestMustMake_CorrectType(t *testing.T) {
	c := New()
	c.Instance("greeting", "hello")

	result := MustMake[string](c, "greeting")
	if result != "hello" {
		t.Errorf("expected 'hello', got '%s'", result)
	}
}

// TC-06: MustMake panics on type mismatch
func TestMustMake_PanicsOnTypeMismatch(t *testing.T) {
	c := New()
	c.Instance("num", 42)

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic on type mismatch")
		}
	}()

	_ = MustMake[string](c, "num")
}

// TC-07: Has returns true for bound service
func TestHas_BoundService(t *testing.T) {
	c := New()
	c.Bind("svc", func(_ *Container) interface{} { return "x" })

	if !c.Has("svc") {
		t.Error("Has should return true for bound service")
	}
}

// TC-08: Has returns true for instance
func TestHas_Instance(t *testing.T) {
	c := New()
	c.Instance("cfg", "value")

	if !c.Has("cfg") {
		t.Error("Has should return true for instance")
	}
}

// TC-09: Has returns false for unregistered service
func TestHas_Unregistered(t *testing.T) {
	c := New()

	if c.Has("nonexistent") {
		t.Error("Has should return false for unregistered service")
	}
}

// TC-10: Instance takes priority over binding in Make
func TestMake_InstancePriorityOverBinding(t *testing.T) {
	c := New()
	c.Bind("svc", func(_ *Container) interface{} { return "from-binding" })
	c.Instance("svc", "from-instance")

	result := c.Make("svc").(string)
	if result != "from-instance" {
		t.Errorf("expected 'from-instance', got '%s'", result)
	}
}

// TC-11: Bind overwrites previous binding
func TestBind_OverwritesPrevious(t *testing.T) {
	c := New()
	c.Bind("svc", func(_ *Container) interface{} { return "A" })
	c.Bind("svc", func(_ *Container) interface{} { return "B" })

	result := c.Make("svc").(string)
	if result != "B" {
		t.Errorf("expected 'B' (last-write-wins), got '%s'", result)
	}
}

// TC-12: Concurrent Make on singleton is safe
func TestSingleton_ConcurrentSafe(t *testing.T) {
	c := New()
	var factoryCalls int64
	c.Singleton("svc", func(_ *Container) interface{} {
		atomic.AddInt64(&factoryCalls, 1)
		return "singleton-value"
	})

	var wg sync.WaitGroup
	results := make([]interface{}, 100)
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			results[idx] = c.Make("svc")
		}(i)
	}
	wg.Wait()

	// All results should be the same value
	for i, r := range results {
		if r != "singleton-value" {
			t.Errorf("goroutine %d got unexpected value: %v", i, r)
		}
	}

	// With double-check locking, factory should be called exactly once
	calls := atomic.LoadInt64(&factoryCalls)
	if calls != 1 {
		t.Errorf("factory should be called exactly once, called %d times", calls)
	}
}

// TC-16: Instance overwrites previous instance with same name
func TestInstance_OverwritesPrevious(t *testing.T) {
	c := New()
	c.Instance("svc", "objA")
	c.Instance("svc", "objB")

	result := c.Make("svc").(string)
	if result != "objB" {
		t.Errorf("expected 'objB' (last-write-wins), got '%s'", result)
	}
}

// TC-17: Bind factory resolves another service from container
func TestBind_FactoryResolvesDependency(t *testing.T) {
	c := New()

	type Config struct{ DSN string }
	type DB struct{ Config *Config }

	c.Instance("config", &Config{DSN: "postgres://localhost"})
	c.Bind("db", func(cont *Container) interface{} {
		cfg := MustMake[*Config](cont, "config")
		return &DB{Config: cfg}
	})

	db := MustMake[*DB](c, "db")
	if db.Config.DSN != "postgres://localhost" {
		t.Errorf("expected 'postgres://localhost', got '%s'", db.Config.DSN)
	}
}

// TC-18: TryMake returns error for missing service (no panic)
func TestTryMake_MissingService(t *testing.T) {
	c := New()
	_, err := c.TryMake("missing")
	if err == nil {
		t.Fatal("TryMake should return error for missing service")
	}
}

// TC-19: TryMake returns value for registered service
func TestTryMake_RegisteredService(t *testing.T) {
	c := New()
	c.Instance("greeting", "hello")
	val, err := c.TryMake("greeting")
	if err != nil {
		t.Fatalf("TryMake error: %v", err)
	}
	if val != "hello" {
		t.Errorf("TryMake = %v, want %q", val, "hello")
	}
}

// TC-20: Generic TryMake catches type assertion failure
func TestTryMakeGeneric_TypeMismatch(t *testing.T) {
	c := New()
	c.Instance("count", 42)
	_, err := TryMake[string](c, "count")
	if err == nil {
		t.Fatal("TryMake[string] should return error for int service")
	}
}

// TC-21: Generic TryMake returns correct type
func TestTryMakeGeneric_CorrectType(t *testing.T) {
	c := New()
	c.Instance("name", "rapidgo")
	val, err := TryMake[string](c, "name")
	if err != nil {
		t.Fatalf("TryMake error: %v", err)
	}
	if val != "rapidgo" {
		t.Errorf("TryMake = %q, want %q", val, "rapidgo")
	}
}
