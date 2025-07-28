// logo/writer_channel.go
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
func NewChannelWriter(ch chan string) *ChannelWriter {
	return &ChannelWriter{ch: ch}
}

// Write implements the io.Writer interface for ChannelWriter.
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
