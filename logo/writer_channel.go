// Package logger provides functionality for structured logging.
//
// This file contains the channel writer implementation which allows log messages
// to be sent to a Go channel for asynchronous processing.
package logger

import (
	"strings"
)

// ChannelWriter is a custom writer that sends log messages to a channel.
// It implements the io.Writer interface, allowing it to be used with standard logging functions.
// This is useful for applications that need to collect log messages asynchronously or process them in a different way.
type ChannelWriter struct {
	ch chan string
}

// NewChannelWriter creates a new ChannelWriter instance.
// It initializes the writer with the provided channel.
// The channel should be buffered to handle log messages without blocking.
//
// Parameters:
//   - ch: A channel of strings that will receive the log messages
//
// Returns a ChannelWriter that implements io.Writer
func NewChannelWriter(ch chan string) *ChannelWriter {
	return &ChannelWriter{ch: ch}
}

// Write implements the io.Writer interface for ChannelWriter.
// It sends the log message to the channel in a non-blocking way.
//
// Parameters:
//   - p: The byte slice containing the log message to write
//
// Returns the number of bytes processed and any error encountered.
// Note that if the channel is full, the message will be dropped but no error is returned.
func (cw *ChannelWriter) Write(p []byte) (int, error) {
	msg := strings.TrimSpace(string(p))
	if msg != "" {
		select {
		case cw.ch <- msg:
			// Successfully sent to channel
		default:
			// Channel is full, drop message
		}
	} else {
		// If the message is empty, we do not send it to the channel
		// This prevents sending empty log messages
		return 0, nil
	}
	return len(p), nil
}
