package logger

import (
	"testing"
	"time"
)

func TestNewChannelWriter(t *testing.T) {
	// Suppress log output for this test
	defer SuppressLogOutput(t)()

	ch := make(chan string, 10)
	writer := NewChannelWriter(ch)

	if writer == nil {
		t.Fatal("NewChannelWriter returned nil")
	}

	if writer.ch != ch {
		t.Error("ChannelWriter's channel doesn't match the provided channel")
	}
}

func TestChannelWriter_Write_SuccessfulWrite(t *testing.T) {
	// Suppress log output for this test
	defer SuppressLogOutput(t)()

	ch := make(chan string, 5)
	writer := NewChannelWriter(ch)

	testMsg := "test message"
	n, err := writer.Write([]byte(testMsg))

	if err != nil {
		t.Errorf("Write() returned error: %v", err)
	}

	if n != len(testMsg) {
		t.Errorf("Write() returned %d bytes written, want %d", n, len(testMsg))
	}

	// Check if message was sent to the channel
	select {
	case msg := <-ch:
		if msg != testMsg {
			t.Errorf("Channel received %q, want %q", msg, testMsg)
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("Timeout waiting for message on channel")
	}
}

func TestChannelWriter_Write_EmptyMessage(t *testing.T) {
	// Suppress log output for this test
	defer SuppressLogOutput(t)()

	ch := make(chan string, 5)
	writer := NewChannelWriter(ch)

	// Write an empty message - should be ignored
	n, err := writer.Write([]byte(""))

	if err != nil {
		t.Errorf("Write() returned error: %v", err)
	}

	if n != 0 {
		t.Errorf("Write() returned %d bytes written for empty message, want 0", n)
	}

	// Check that nothing was sent to the channel
	select {
	case msg := <-ch:
		t.Errorf("Channel received unexpected message %q", msg)
	case <-time.After(100 * time.Millisecond):
		// This is the expected behavior - no message should be received
	}
}

func TestChannelWriter_Write_WhitespaceMessage(t *testing.T) {
	// Suppress log output for this test
	defer SuppressLogOutput(t)()

	ch := make(chan string, 5)
	writer := NewChannelWriter(ch)

	// Write a whitespace-only message - should be ignored
	n, err := writer.Write([]byte("  \n\t  "))

	if err != nil {
		t.Errorf("Write() returned error: %v", err)
	}

	// The implementation returns the original bytes written, even for whitespace messages
	expectedBytes := len("  \n\t  ")
	if n != expectedBytes {
		t.Errorf("Write() returned %d bytes written, want %d", n, expectedBytes)
	}

	// Check that nothing was sent to the channel
	select {
	case msg := <-ch:
		t.Errorf("Channel received unexpected message %q", msg)
	case <-time.After(100 * time.Millisecond):
		// This is the expected behavior - no message should be received
	}
}

func TestChannelWriter_Write_FullChannel(t *testing.T) {
	// Suppress log output for this test
	defer SuppressLogOutput(t)()

	// Create a channel with capacity 1
	ch := make(chan string, 1)
	writer := NewChannelWriter(ch)

	// Fill the channel
	ch <- "existing message"

	// Try to write - should not block and should not send to channel
	testMsg := "test message that won't be sent"
	n, err := writer.Write([]byte(testMsg))

	if err != nil {
		t.Errorf("Write() returned error: %v", err)
	}

	if n != len(testMsg) {
		t.Errorf("Write() returned %d bytes written, want %d", n, len(testMsg))
	}

	// Check that the channel still contains only the original message
	select {
	case msg := <-ch:
		if msg != "existing message" {
			t.Errorf("Channel contained %q, want %q", msg, "existing message")
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("Timeout reading from channel")
	}

	// Verify no additional messages were added (channel should be empty now)
	select {
	case msg := <-ch:
		t.Errorf("Channel unexpectedly contained extra message: %q", msg)
	case <-time.After(100 * time.Millisecond):
		// This is the expected behavior - no additional messages
	}
}

func TestChannelWriter_Write_MultipleMessages(t *testing.T) {
	// Suppress log output for this test
	defer SuppressLogOutput(t)()

	ch := make(chan string, 5)
	writer := NewChannelWriter(ch)

	messages := []string{
		"first message",
		"second message",
		"third message",
	}

	// Write multiple messages
	for _, msg := range messages {
		n, err := writer.Write([]byte(msg))

		if err != nil {
			t.Errorf("Write() returned error: %v", err)
		}

		if n != len(msg) {
			t.Errorf("Write() returned %d bytes written, want %d", n, len(msg))
		}
	}

	// Verify all messages were sent to the channel in order
	for i, expectedMsg := range messages {
		select {
		case msg := <-ch:
			if msg != expectedMsg {
				t.Errorf("Message %d: got %q, want %q", i, msg, expectedMsg)
			}
		case <-time.After(100 * time.Millisecond):
			t.Errorf("Timeout waiting for message %d", i)
		}
	}

	// Verify no extra messages
	select {
	case msg := <-ch:
		t.Errorf("Channel unexpectedly contained extra message: %q", msg)
	case <-time.After(100 * time.Millisecond):
		// Expected - no more messages
	}
}
