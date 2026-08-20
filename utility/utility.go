package utility

import (
	"path/filepath"
	"strings"
)

func isSupportedAsset(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))

	switch ext {
	case ".jpg", ".jpeg", ".png", ".webp", ".gif":
		return true
	default:
		return false
	}
}
