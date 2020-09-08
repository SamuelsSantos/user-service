package config

import (
	"os"
	"strings"
)

// GetenvString get enviroment value or default value
func GetenvString(key, standard string) string {
	value := os.Getenv(key)
	if len(strings.TrimSpace(value)) == 0 {
		return standard
	}
	return value
}
