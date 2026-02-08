package logo

import (
	"fmt"
	"strings"
)

// ReservedKeys are field names reserved by the logging schema.
// These keys cannot be overwritten by user-provided fields.
var ReservedKeys = map[string]bool{
	"ts":     true, // timestamp
	"level":  true, // log level
	"msg":    true, // message
	"logger": true, // logger name
	"source": true, // source location
	"caller": true, // caller info (alternative name)
}

// IsReservedKey returns true if the key is reserved by the schema.
func IsReservedKey(key string) bool {
	return ReservedKeys[key]
}

// ValidateFields checks for reserved key conflicts and returns an error if any are found.
func ValidateFields(fields map[string]interface{}) error {
	var conflicts []string
	for key := range fields {
		if IsReservedKey(key) {
			conflicts = append(conflicts, key)
		}
	}

	if len(conflicts) > 0 {
		return fmt.Errorf("fields contain reserved keys: %s", strings.Join(conflicts, ", "))
	}

	return nil
}

// FilterReservedKeys returns a new map with reserved keys removed.
// This is a safe operation that prevents schema violations.
func FilterReservedKeys(fields map[string]interface{}) map[string]interface{} {
	filtered := make(map[string]interface{}, len(fields))
	for key, value := range fields {
		if !IsReservedKey(key) {
			filtered[key] = value
		}
	}
	return filtered
}

// Schema field names - these are the canonical names used in JSON output
const (
	FieldTimestamp = "ts"
	FieldLevel     = "level"
	FieldMessage   = "msg"
	FieldLogger    = "logger"
	FieldSource    = "source"
)
