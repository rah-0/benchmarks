package util

import (
	"os"
)

func GetEnvVariable(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return ""
	}
	return value
}
