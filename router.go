package logo

import (
	"sync"
)

// Router dispatches log entries to multiple handlers.
// It implements a fan-out pattern where each entry is sent to all registered handlers.
type Router struct {
	handlers []Handler
	mu       sync.RWMutex
}

// NewRouter creates a new router with the given handlers.
func NewRouter(handlers ...Handler) *Router {
	return &Router{
		handlers: handlers,
	}
}

// Route sends the log entry to all registered handlers.
func (r *Router) Route(entry *Entry) {
	r.mu.RLock()
	handlers := r.handlers
	r.mu.RUnlock()

	// Fan out to all handlers
	for _, handler := range handlers {
		// Each handler decides whether to process the entry based on its level
		handler.Handle(entry)
	}
}

// AddHandler adds a new handler to the router.
func (r *Router) AddHandler(handler Handler) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.handlers = append(r.handlers, handler)
}

// Handlers returns a copy of the current handlers.
func (r *Router) Handlers() []Handler {
	r.mu.RLock()
	defer r.mu.RUnlock()

	handlers := make([]Handler, len(r.handlers))
	copy(handlers, r.handlers)
	return handlers
}

// Close closes all handlers.
func (r *Router) Close() error {
	r.mu.RLock()
	handlers := r.handlers
	r.mu.RUnlock()

	var lastErr error
	for _, handler := range handlers {
		if err := handler.Close(); err != nil {
			lastErr = err
		}
	}
	return lastErr
}
