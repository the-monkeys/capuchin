package config

import (
	"fmt"
	"strings"
)

func validateDatabaseConfig(cfg AppConfig) error {
	missing := make([]string, 0, 4)

	if cfg.PostgresUser == "" {
		missing = append(missing, "POSTGRES_USER")
	}
	if cfg.PostgresPassword == "" {
		missing = append(missing, "POSTGRES_PASSWORD")
	}
	if cfg.PostgresDB == "" {
		missing = append(missing, "POSTGRES_DB")
	}
	if cfg.PostgresHost == "" {
		missing = append(missing, "POSTGRES_HOST")
	}
	if len(missing) > 0 {
		return fmt.Errorf("missing required database config: %s", strings.Join(missing, ", "))
	}

	return nil
}
