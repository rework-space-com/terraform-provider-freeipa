package freeipa

import (
	"os"
	"strconv"
	"strings"
)

func getEnvAsBool(name string, defaultVal bool) bool {
	valStr := os.Getenv(name)
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}

	return defaultVal
}

	return defaultVal
}

func decodeSlash(str string) string {
	return strings.ReplaceAll(str, "%2F", string('/'))
}
