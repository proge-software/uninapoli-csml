package envext

import "os"

// GetEnvOrDefault returns the value of the env variable with
// name `key` if present otherwise `default`
func GetEnvOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
