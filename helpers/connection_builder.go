package helpers

import (
	"fmt"
	appConfig "go-cache-me/configs"
)

// ConnectionURLBuilder func for building URL connection.
func ConnectionURLBuilder(n string) (string, error) {
	// Define URL to connection.
	var url string

	// Switch given names.
	switch n {
	case "redis":
		// TODO: Implement Redis connection.
	case "fiber":
		// URL for Fiber connection.
		url = fmt.Sprintf(
			"%s:%s",
			appConfig.GetEnv("SERVER_HOST"),
			appConfig.GetEnv("SERVER_PORT"),
		)
	default:
		// Return error message.
		return "", fmt.Errorf("connection name '%v' is not supported", n)
	}

	// Return connection URL.
	return url, nil
}
