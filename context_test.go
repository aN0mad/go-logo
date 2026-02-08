package logo

import (
	"context"
	"testing"
)

func TestWithLogger(t *testing.T) {
	logger := New(WithName("test"))
	ctx := WithLogger(context.Background(), logger)

	retrieved := FromContext(ctx)
	if retrieved.name != "test" {
		t.Errorf("expected logger name 'test', got '%s'", retrieved.name)
	}
}

func TestWithLogger_NilContext(t *testing.T) {
	logger := New(WithName("test"))
	ctx := WithLogger(nil, logger)

	if ctx == nil {
		t.Error("expected non-nil context")
	}

	retrieved := FromContext(ctx)
	if retrieved.name != "test" {
		t.Errorf("expected logger name 'test', got '%s'", retrieved.name)
	}
}

func TestFromContext_NoLogger(t *testing.T) {
	ctx := context.Background()
	logger := FromContext(ctx)

	if logger == nil {
		t.Error("expected default logger, got nil")
	}

	// Should return a valid logger that can be used
	logger.Info("test message")
}

func TestFromContext_NilContext(t *testing.T) {
	logger := FromContext(nil)

	if logger == nil {
		t.Error("expected default logger, got nil")
	}
}

func TestWithRequestID(t *testing.T) {
	ctx := WithRequestID(context.Background(), "req-123")

	requestID := RequestIDFromContext(ctx)
	if requestID != "req-123" {
		t.Errorf("expected request ID 'req-123', got '%s'", requestID)
	}
}

func TestWithRequestID_NilContext(t *testing.T) {
	ctx := WithRequestID(nil, "req-456")

	if ctx == nil {
		t.Error("expected non-nil context")
	}

	requestID := RequestIDFromContext(ctx)
	if requestID != "req-456" {
		t.Errorf("expected request ID 'req-456', got '%s'", requestID)
	}
}

func TestRequestIDFromContext_NoID(t *testing.T) {
	ctx := context.Background()
	requestID := RequestIDFromContext(ctx)

	if requestID != "" {
		t.Errorf("expected empty request ID, got '%s'", requestID)
	}
}

func TestRequestIDFromContext_NilContext(t *testing.T) {
	requestID := RequestIDFromContext(nil)

	if requestID != "" {
		t.Errorf("expected empty request ID, got '%s'", requestID)
	}
}

func TestContextWithFields(t *testing.T) {
	logger := New(WithName("base"))
	ctx := WithLogger(context.Background(), logger)

	newLogger := ContextWithFields(ctx, F("key", "value"))

	// Verify the new logger has the field
	if len(newLogger.fields) != 1 {
		t.Errorf("expected 1 field, got %d", len(newLogger.fields))
	}

	if newLogger.fields["key"] != "value" {
		t.Errorf("expected field value 'value', got '%v'", newLogger.fields["key"])
	}

	// Original logger should be unchanged
	if len(logger.fields) != 0 {
		t.Error("original logger should not be modified")
	}
}

func TestL_Shorthand(t *testing.T) {
	logger := New(WithName("shorthand"))
	ctx := WithLogger(context.Background(), logger)

	retrieved := L(ctx)
	if retrieved.name != "shorthand" {
		t.Errorf("expected logger name 'shorthand', got '%s'", retrieved.name)
	}
}

func TestContextChaining(t *testing.T) {
	// Test that we can chain context operations
	ctx := context.Background()

	logger := New(WithName("chained"))
	ctx = WithLogger(ctx, logger)
	ctx = WithRequestID(ctx, "req-789")

	// Verify both values are present
	retrievedLogger := FromContext(ctx)
	if retrievedLogger.name != "chained" {
		t.Errorf("expected logger name 'chained', got '%s'", retrievedLogger.name)
	}

	retrievedID := RequestIDFromContext(ctx)
	if retrievedID != "req-789" {
		t.Errorf("expected request ID 'req-789', got '%s'", retrievedID)
	}
}

func TestContextIsolation(t *testing.T) {
	// Verify that contexts don't interfere with each other
	ctx1 := WithRequestID(context.Background(), "req-1")
	ctx2 := WithRequestID(context.Background(), "req-2")

	id1 := RequestIDFromContext(ctx1)
	id2 := RequestIDFromContext(ctx2)

	if id1 != "req-1" {
		t.Errorf("context 1: expected 'req-1', got '%s'", id1)
	}
	if id2 != "req-2" {
		t.Errorf("context 2: expected 'req-2', got '%s'", id2)
	}
}
