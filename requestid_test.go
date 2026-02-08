package logo

import (
	"strings"
	"testing"
)

func TestUUIDGenerator(t *testing.T) {
	gen := &UUIDGenerator{}

	// Generate multiple IDs
	ids := make(map[string]bool)
	for i := 0; i < 100; i++ {
		id := gen.Generate()

		// Check format (UUID v4 format: xxxxxxxx-xxxx-4xxx-xxxx-xxxxxxxxxxxx)
		if len(id) != 36 {
			t.Errorf("expected UUID length 36, got %d: %s", len(id), id)
		}

		// Check for duplicates
		if ids[id] {
			t.Errorf("duplicate ID generated: %s", id)
		}
		ids[id] = true

		// Verify format
		parts := strings.Split(id, "-")
		if len(parts) != 5 {
			t.Errorf("expected 5 UUID parts, got %d: %s", len(parts), id)
		}
	}
}

func TestShortIDGenerator(t *testing.T) {
	gen := &ShortIDGenerator{}

	// Generate multiple IDs
	ids := make(map[string]bool)
	for i := 0; i < 100; i++ {
		id := gen.Generate()

		// Check length (16 hex characters)
		if len(id) != 16 {
			t.Errorf("expected short ID length 16, got %d: %s", len(id), id)
		}

		// Check for duplicates
		if ids[id] {
			t.Errorf("duplicate ID generated: %s", id)
		}
		ids[id] = true

		// Verify it's hex
		for _, c := range id {
			if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f')) {
				t.Errorf("expected hex character, got '%c' in: %s", c, id)
			}
		}
	}
}

func TestSequentialGenerator(t *testing.T) {
	gen := NewSequentialGenerator("test")

	// Generate IDs and verify sequence
	for i := 1; i <= 10; i++ {
		id := gen.Generate()
		if !strings.HasPrefix(id, "test-") {
			t.Errorf("expected prefix 'test-', got: %s", id)
		}
	}
}

func TestSequentialGenerator_NoPrefix(t *testing.T) {
	gen := NewSequentialGenerator("")

	id := gen.Generate()

	// Should be numeric only
	for _, c := range id {
		if c < '0' || c > '9' {
			t.Errorf("expected numeric ID, got: %s", id)
		}
	}
}

func TestGenerateRequestID(t *testing.T) {
	// Test default generator
	id := GenerateRequestID()

	if id == "" {
		t.Error("expected non-empty request ID")
	}

	// Should be UUID format by default
	if len(id) != 36 {
		t.Errorf("expected UUID length 36, got %d: %s", len(id), id)
	}
}

func TestSetRequestIDGenerator(t *testing.T) {
	// Save original
	original := DefaultRequestIDGenerator
	defer func() {
		DefaultRequestIDGenerator = original
	}()

	// Set to short ID generator
	SetRequestIDGenerator(&ShortIDGenerator{})

	id := GenerateRequestID()

	// Should be short format
	if len(id) != 16 {
		t.Errorf("expected short ID length 16, got %d: %s", len(id), id)
	}
}

func TestRequestIDUniqueness(t *testing.T) {
	// Generate a large number of IDs and check for duplicates
	ids := make(map[string]bool)
	count := 10000

	for i := 0; i < count; i++ {
		id := GenerateRequestID()
		if ids[id] {
			t.Errorf("duplicate ID found: %s", id)
		}
		ids[id] = true
	}

	if len(ids) != count {
		t.Errorf("expected %d unique IDs, got %d", count, len(ids))
	}
}

func BenchmarkUUIDGenerator(b *testing.B) {
	gen := &UUIDGenerator{}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = gen.Generate()
	}
}

func BenchmarkShortIDGenerator(b *testing.B) {
	gen := &ShortIDGenerator{}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = gen.Generate()
	}
}

func BenchmarkSequentialGenerator(b *testing.B) {
	gen := NewSequentialGenerator("bench")
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = gen.Generate()
	}
}
