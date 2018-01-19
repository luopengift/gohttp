package gohttp

import (
	"strings"
)

// If s starts with one of suffixs; return ture
func hasSuffixs(s string, suffixs ...string) bool {
	for _, suffix := range suffixs {
		if ok := strings.HasSuffix(s, suffix); ok {
			return true
		}
	}
	return false
}

func hasPrefixs(s string, prefixs ...string) bool {
	for _, prefix := range prefixs {
		if ok := strings.HasPrefix(s, prefix); ok {
			return true
		}
	}
	return false
}
