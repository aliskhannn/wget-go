package parser

import (
	"bytes"
	"net/url"
	"path/filepath"
	"strings"
)

// NormalizeURL ensures the URL has a scheme (http/https).
// If missing, defaults to https.
func NormalizeURL(raw string) (*url.URL, error) {
	raw = strings.TrimSpace(raw)

	if !strings.HasPrefix(raw, "https://") {
		raw = "https://" + raw
	}

	return url.Parse(raw)
}

// IsHTML reports whether the given file path and data most likely represent
// an HTML document. It checks the file extension first, and if uncertain,
// falls back to detecting a "<!DOCTYPE html>" prefix in the data.
func IsHTML(path string, data []byte) bool {
	// Check common HTML file extensions.
	ext := filepath.Ext(path)
	if ext == ".html" || ext == ".htm" || ext == "" {
		return true
	}

	// Fallback: detect HTML signature in the content itself.
	return bytes.HasPrefix(data, []byte("<!DOCTYPE html>"))
}
