package logger

import (
	"os"

	"github.com/rs/zerolog"
)

// Log is a shared application logger.
var Log zerolog.Logger

// Init initializes the global logger.
func Init() {
	Log = zerolog.New(os.Stdout).With().Timestamp().Logger()
}

