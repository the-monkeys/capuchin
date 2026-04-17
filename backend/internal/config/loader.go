package config

import (
	"fmt"
	"os"
	"strconv"
)

func postgresPortFromEnv() (int, error) {
	rawPort := os.Getenv("POSTGRES_PORT")
	if rawPort == "" {
		return 5432, nil
	}

	port, err := strconv.Atoi(rawPort)
	if err != nil {
		return 0, fmt.Errorf("invalid POSTGRES_PORT %q: must be numeric", rawPort)
	}
	if port <= 0 {
		return 0, fmt.Errorf("invalid POSTGRES_PORT %q: must be greater than 0", rawPort)
	}

	return port, nil
}
