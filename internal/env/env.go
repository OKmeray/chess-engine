package env

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// parseEnv parses the environment variable for key, returning fallback if unset or an error if malformed.
func parseEnv[T any](key string, fallback T, parser func(string) (T, error)) (T, error) {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback, nil
	}

	parsed, err := parser(val)
	if err != nil {
		var zero T
		return zero, fmt.Errorf("env: invalid value %q for key %q: %w", val, key, err)
	}

	return parsed, nil
}

// GetString returns the string environment variable for key, or fallback if unset.
func GetString(key, fallback string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	return val
}

// GetInt returns the integer environment variable for key, or fallback if unset.
func GetInt(key string, fallback int) (int, error) {
	return parseEnv(key, fallback, strconv.Atoi)
}

// GetDuration returns the time.Duration environment variable for key, or fallback if unset.
func GetDuration(key string, fallback time.Duration) (time.Duration, error) {
	return parseEnv(key, fallback, time.ParseDuration)
}
