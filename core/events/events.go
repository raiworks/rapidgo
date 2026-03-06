package events

import "sync"

// Handler is a function that handles an event.
type Handler func(payload interface{})

// Dispatcher manages event listeners and dispatches events.
type Dispatcher struct {
	mu        sync.RWMutex
	listeners map[string][]Handler
}

// NewDispatcher returns an empty event dispatcher.
func NewDispatcher() *Dispatcher {
	return &Dispatcher{listeners: make(map[string][]Handler)}
}

// Listen registers a handler for the named event.
func (d *Dispatcher) Listen(event string, handler Handler) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.listeners[event] = append(d.listeners[event], handler)
}

// Dispatch fires all handlers for the event asynchronously.
func (d *Dispatcher) Dispatch(event string, payload interface{}) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	for _, h := range d.listeners[event] {
		go h(payload)
	}
}

// DispatchSync fires all handlers for the event sequentially.
func (d *Dispatcher) DispatchSync(event string, payload interface{}) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	for _, h := range d.listeners[event] {
		h(payload)
	}
}

// Has returns true if any listener is registered for the event.
func (d *Dispatcher) Has(event string) bool {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return len(d.listeners[event]) > 0
}
