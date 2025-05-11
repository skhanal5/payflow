package config

import (
	"fmt"
	"os"
)

func GetEnvOrPanic(key string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	panic(fmt.Sprintf("required environment variable %q not set", key)) //TODO: Use logger
}