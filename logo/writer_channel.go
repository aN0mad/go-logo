// Package logo provides functionality for structured logging.
//
// This file contains the channel writer implementation which allows log messages
// to be sent to a Go channel for asynchronous processing.
package logo

import (
	"strings"
)

// ChannelWriter is an io.Writer that sends log messages to a channel.
// It can be used to pass log messages to custom processing routines.
type ChannelWriter struct {
	ch chan string
}

// NewChannelWriter creates a new writer that sends log messages to the given channel.
//
// Parameters:
//   - ch: The channel to which log messages will be sent
//
// Returns:
//   - *ChannelWriter: A channel writer that implements io.Writer
func NewChannelWriter(ch chan string) *ChannelWriter {
	return &ChannelWriter{ch: ch}
}

// Write implements the io.Writer interface for ChannelWriter.
// It sends the log message to the channel associated with this writer.
//
// Parameters:
//   - p: The byte slice containing the log message
//
// Returns:
//   - int: The number of bytes processed
//   - error: Any error encountered during writing (or nil if successful)
func (cw *ChannelWriter) Write(p []byte) (int, error) {
	msg := strings.TrimSpace(string(p))
	if msg != "" {
		select {
		case cw.ch <- msg:
			// Successfully sent to channel
		default:
			// Channel is full, drop message
		}
	}
	// Always return original length, even if message is filtered
	return len(p), nil
}
