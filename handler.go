package logo

// Handler processes log entries and outputs them to a destination.
// Handlers define the minimum log level, formatting, and output destination.
type Handler interface {
	// Handle processes a log entry.
	// Returns true if the entry was handled, false if it was filtered.
	Handle(entry *Entry) bool

	// Level returns the minimum level this handler will process.
	Level() Level

	// Enabled returns whether this handler will process the given level.
	Enabled(level Level) bool

	// Close releases any resources held by the handler.
	Close() error
}

// BaseHandler provides common handler functionality.
type BaseHandler struct {
	minLevel Level
}

// NewBaseHandler creates a new base handler with the specified minimum level.
func NewBaseHandler(level Level) *BaseHandler {
	return &BaseHandler{
		minLevel: level,
	}
}

// Level returns the minimum level this handler will process.
func (h *BaseHandler) Level() Level {
	return h.minLevel
}

// Enabled returns whether this handler will process the given level.
func (h *BaseHandler) Enabled(level Level) bool {
	return level >= h.minLevel
}

// Close is a no-op for the base handler.
func (h *BaseHandler) Close() error {
	return nil
}
