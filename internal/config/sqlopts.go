package config

import (
	"fmt"
)

func SqlOpts(c DatabaseConfig) string {
	base := fmt.Sprintf("%s:%s/%s", c.Host, c.Port, c.DB)
	if c.SSLMode != "" {
		base = fmt.Sprintf("%s?sslmode=%s", base, c.SSLMode)
	}

	userpass := ""
	if c.User != "" {
		userpass = c.User
	}

	if c.Pass != "" {
		userpass = fmt.Sprintf("%s:%s", userpass, c.Pass)
	}

	if userpass != "" {
		base = fmt.Sprintf("%s@%s", userpass, base)
	}

	return fmt.Sprintf("postgres://%s", base)
}
