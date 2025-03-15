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

// Some resource names are use construct the resource Id (name/cat/value). If they contain a slash, it messes up with the parsing of resources id.
func encodeSlash(str string) string {
	return strings.ReplaceAll(str, string('/'), "%2F")
}

func decodeSlash(str string) string {
	return strings.ReplaceAll(str, "%2F", string('/'))
}

func isStringListContainsCaseInsensistive(strList *[]string, str *string) bool {
	for _, s := range *strList {
		if strings.EqualFold(s, *str) {
			return true
		}
	}
	return false
}
