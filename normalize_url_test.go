package main

import (
	"strings"
	"testing"
)

func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		name     string
		inputURL string
		expected string
	}{
		{
			name:     "remove https scheme",
			inputURL: "https://blog.boot.dev/path",
			expected: "blog.boot.dev/path",
		},
		{
			name:     "remove https schema and end slash",
			inputURL: "https://blog.boot.dev/path/",
			expected: "blog.boot.dev/path",
		},
		{
			name:     "remove http schema",
			inputURL: "http://blog.boot.dev/path",
			expected: "blog.boot.dev/path",
		},
		{
			name:     "remove http schema and end slash",
			inputURL: "http://blog.boot.dev/path/",
			expected: "blog.boot.dev/path",
		},
		{
			name:     "no shema remove end slash",
			inputURL: "blog.boot.dev/path/",
			expected: "blog.boot.dev/path",
		},
		{
			name:     "remove shema and change to lowercase remove end slash",
			inputURL: "http://blog.BOOT.dev/path",
			expected: "blog.boot.dev/path",
		},
		{
			name:     "improper shema remove end slash",
			inputURL: "://blog.boot.dev/path/",
			expected: "missing protocol scheme",
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := normalizeURL(tc.inputURL)
			if err != nil {
				if strings.Contains(err.Error(), tc.expected) {
					return
				}
				t.Errorf("Test %v - '%s' FAIL: unexpected error: %v", i, tc.name, err)
				return
			}
			if actual != tc.expected {
				t.Errorf("Test %v - %s FAIL: expected URL: %v, actual: %v", i, tc.name, tc.expected, actual)
			}
		})
	}
}
