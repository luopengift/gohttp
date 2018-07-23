package gohttp

import (
	"net/http"
	"os"
	"path/filepath"
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

func toHTTPError(err error) (string, int) {
	if os.IsNotExist(err) {
		return http.StatusText(http.StatusNotFound), http.StatusNotFound
	}
	if os.IsPermission(err) {
		return http.StatusText(http.StatusForbidden), http.StatusForbidden
	}
	// Default:
	return http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError
}

// WalkDir walk dir path
func WalkDir(path string, matchs ...func(path string) bool) ([]string, error) {
	match := func(string) bool { return true }
	if len(matchs) > 0 && matchs[0] != nil {
		match = matchs[0]
	}

	var files []string
	err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			return nil
		}
		if match(path) {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}
