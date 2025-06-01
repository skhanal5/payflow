package config

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

func InitLogger(environment string) zerolog.Logger {
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.TimestampFieldName = "ts"
	zerolog.LevelFieldName = "level"

	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}

	// Base logger with common context
	logger := zerolog.New(output).With().
		Timestamp().
		Str("service", "inventory").
		Str("env", environment).
		Logger()

	return logger
}
