package utils

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
)

// DecodeImage reads an image file from disk and decodes it.
// Supports PNG, JPEG, GIF, BMP (via standard library).
func DecodeImage(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image (supported formats: PNG, JPEG, GIF, BMP): %w", err)
	}
	return img, nil
}

// ExpandWildcard expands a pattern (e.g., "*.png") to a list of matching files.
// If pattern contains no wildcards, returns the pattern as a single-element list.
func ExpandWildcard(pattern string) ([]string, error) {
	// Check if pattern contains any wildcard characters
	if !containsWildcard(pattern) {
		// If the file exists, return it; otherwise assume it's a literal path
		if _, err := os.Stat(pattern); err == nil {
			return []string{pattern}, nil
		}
		// Maybe it's a glob that didn't match? Let's treat as literal.
		return []string{pattern}, nil
	}

	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}
	if len(matches) == 0 {
		return nil, fmt.Errorf("no files match pattern: %s", pattern)
	}
	return matches, nil
}

func containsWildcard(s string) bool {
	for _, ch := range s {
		if ch == '*' || ch == '?' {
			return true
		}
	}
	return false
}
