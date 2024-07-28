package freeipa

import (
	"strings"
)

func utilsGetArry(itemsRaw []interface{}) []string {
	res := make([]string, len(itemsRaw))
	for i, raw := range itemsRaw {
		res[i] = raw.(string)
	}
	return res
}

// Some resource names are use construct the resource Id (name/cat/value). If they contain a slash, it messes up with the parsing of resources id.
func encodeSlash(str string) string {
	return strings.ReplaceAll(str, string('/'), "%2F")
}

func decodeSlash(str string) string {
	return strings.ReplaceAll(str, "%2F", string('/'))
}