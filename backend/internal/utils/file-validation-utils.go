package utils

import (
	"archive/zip"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gianghp/statify/internal/core"
)

func ValidateZipArchive(files []*zip.File) error {
	var totalUncompressedSize int64

	if len(files) > core.MaxFileCount {
		return fmt.Errorf("Archive contains too many files: limit is %d", core.MaxFileCount)
	}

	for _, file := range files {
		// 1. Path Traversal Protection
		// Ensures file names don't look like "../../etc/passwd"
		if strings.Contains(file.Name, "..") || strings.HasPrefix(file.Name, "/") {
			return fmt.Errorf("Invalid file path in zip: %s", file.Name)
		}

		// 2. Resource Limit Check (Disk Fill Prevention)
		totalUncompressedSize += int64(file.UncompressedSize64)
		if totalUncompressedSize > core.MaxDeploymentSize {
			return fmt.Errorf("Uncompressed size exceeds limit of %d bytes", core.MaxDeploymentSize)
		}
	}

	return nil
}

var AllowedStaticExtensions = map[string]bool{
	// Structure & Logic
	".html": true, ".htm": true, ".xml": true,
	".css": true,
	".js":  true, ".mjs": true, ".json": true, ".map": true,

	// Images
	".png": true, ".jpg": true, ".jpeg": true, ".gif": true,
	".svg": true, ".ico": true, ".webp": true, ".avif": true,

	// Fonts
	".woff": true, ".woff2": true, ".ttf": true, ".eot": true, ".otf": true,

	// Documents/Text
	".txt": true, ".pdf": true, ".md": true, ".csv": true,

	// Media (Optional - remove if you want to save bandwidth)
	".mp4": true, ".webm": true, ".mp3": true, ".wav": true,
}

func ValidateStaticFileTypes(files []*zip.File) error {
	for _, file := range files {
		if file.FileInfo().IsDir() {
			continue
		}

		fileName := filepath.Base(file.Name)
		if strings.HasPrefix(fileName, ".") || fileName == "__MACOSX" || fileName == "Thumbs.db" {
			continue
		}

		ext := strings.ToLower(filepath.Ext(fileName))
		if !AllowedStaticExtensions[ext] {
			return fmt.Errorf("File type not allowed: %s (in file: %s)", ext, file.Name)
		}
	}
	return nil
}

func ValidateEntrypoint(files []*zip.File) error {
	for _, file := range files {
		if filepath.Clean(file.Name) == "index.html" {
			return nil
		}
	}
	return fmt.Errorf("Missing 'index.html' at the root of the zip file. Please zip the *contents* of your build folder, not the folder itself")
}
