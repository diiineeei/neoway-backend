package env

import (
	"os"
	"strconv"
	"strings"
	"time"
)

// GetString retorna o valor da variável de ambiente ou o valor padrão.
func GetString(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// GetInt retorna o valor da variável de ambiente como int ou o valor padrão.
func GetInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return intValue
}

// GetBool retorna o valor da variável de ambiente como bool ou o valor padrão.
func GetBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}
	return boolValue
}

// GetDuration retorna o valor da variável de ambiente como time.Duration ou o valor padrão.
func GetDuration(key string, defaultValue time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	duration, err := time.ParseDuration(value)
	if err != nil {
		return defaultValue
	}
	return duration
}

// GetStringSlice retorna o valor da variável de ambiente como slice de strings.
// O separador padrão é vírgula (,).
func GetStringSlice(key string, defaultValue []string, separator string) []string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	if separator == "" {
		separator = ","
	}

	parts := strings.Split(value, separator)
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	if len(result) == 0 {
		return defaultValue
	}

	return result
}

// MustGetString retorna o valor da variável de ambiente ou entra em pânico se não existir.
func MustGetString(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic("missing required environment variable: " + key)
	}
	return value
}
