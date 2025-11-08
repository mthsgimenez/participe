package env

import (
	"fmt"
	"os"
	"strconv"
)

func GetString(key string) (string, error) {
	val, ok := os.LookupEnv(key)
	if !ok {
		return "", fmt.Errorf("environment variable %s is not set", key)
	}

	return val, nil
}

func GetInt(key string) (int, error) {
	val, ok := os.LookupEnv(key)
	if !ok {
		return 0, fmt.Errorf("environment variable %s is not set", key)
	}

	intVal, err := strconv.Atoi(val)
	if err != nil {
		return 0, fmt.Errorf("environment variable %s is not an int", key)
	}

	return intVal, nil
}

func GetStringFallback(key, fallback string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	return val
}

func GetIntFallback(key string, fallback int) int {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	intVal, err := strconv.Atoi(val)
	if err != nil {
		return fallback
	}

	return intVal
}
