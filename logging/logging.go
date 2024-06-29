// The package for instantiating a JSON logger
package logging

import (
	"log/slog"
	"os"
)

// GetLogger instantiates the logger for each endpoint without needing to
// repeat oneself.
func GetLogger() *slog.Logger {
	handler := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug})
	logger := slog.New(handler)
	return logger
}
