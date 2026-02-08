package logo

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync/atomic"
	"time"
)

var requestCounter uint64

// RequestIDGenerator defines the interface for generating request IDs.
type RequestIDGenerator interface {
	Generate() string
}

// UUIDGenerator generates UUIDs for request IDs.
type UUIDGenerator struct{}

// Generate creates a new UUID v4.
func (g *UUIDGenerator) Generate() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		// Fallback to timestamp-based ID if random fails
		return fmt.Sprintf("%d-%d", time.Now().UnixNano(), atomic.AddUint64(&requestCounter, 1))
	}

	// Set version (4) and variant bits
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80

	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}

// ShortIDGenerator generates short hex-based request IDs.
// Shorter than UUIDs but still unique enough for most use cases.
type ShortIDGenerator struct{}

// Generate creates a short hex ID (16 characters).
func (g *ShortIDGenerator) Generate() string {
	b := make([]byte, 8)
	_, err := rand.Read(b)
	if err != nil {
		// Fallback to timestamp-based ID
		return fmt.Sprintf("%016x", time.Now().UnixNano())
	}
	return hex.EncodeToString(b)
}

// SequentialGenerator generates sequential request IDs.
// Useful for testing and debugging, not recommended for production.
type SequentialGenerator struct {
	prefix string
}

// NewSequentialGenerator creates a new sequential ID generator with an optional prefix.
func NewSequentialGenerator(prefix string) *SequentialGenerator {
	return &SequentialGenerator{prefix: prefix}
}

// Generate creates a sequential ID.
func (g *SequentialGenerator) Generate() string {
	count := atomic.AddUint64(&requestCounter, 1)
	if g.prefix != "" {
		return fmt.Sprintf("%s-%d", g.prefix, count)
	}
	return fmt.Sprintf("%d", count)
}

// DefaultRequestIDGenerator is the default generator used by middleware.
var DefaultRequestIDGenerator RequestIDGenerator = &UUIDGenerator{}

// GenerateRequestID generates a new request ID using the default generator.
func GenerateRequestID() string {
	return DefaultRequestIDGenerator.Generate()
}

// SetRequestIDGenerator sets the default request ID generator.
func SetRequestIDGenerator(generator RequestIDGenerator) {
	DefaultRequestIDGenerator = generator
}
