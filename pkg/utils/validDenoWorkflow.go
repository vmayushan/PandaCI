package utils

import (
	"path/filepath"
	"slices"
)

func IsValidDenoWorkflow(path string) bool {
	extension := filepath.Ext(path)
	validTypes := []string{".ts", ".js", ".jsx", ".tsx", ".mjs", ".cjs"}

	return slices.Contains(validTypes, extension)
}
