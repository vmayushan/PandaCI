package utils

import (
	"testing"
)

func TestIsValidDenoWorkflow(t *testing.T) {
	tests := []struct {
		path     string
		expected bool
	}{
		{"file.ts", true},
		{"file.js", true},
		{"file.jsx", true},
		{"file.tsx", true},
		{"file.mjs", true},
		{"file.cjs", true},
		{"file.go", false},
		{"file.py", false},
		{"file.txt", false},
		{"file", false},
	}

	for _, test := range tests {
		result := IsValidDenoWorkflow(test.path)
		if result != test.expected {
			t.Errorf("IsValidDenoWorkflow(%s) = %v; want %v", test.path, result, test.expected)
		}
	}
}
