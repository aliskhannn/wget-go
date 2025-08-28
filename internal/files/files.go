package files

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// SaveFile writes content from data to a file derived from the URL inside baseDir.
// It returns the final path where the file was written.
func SaveFile(u *url.URL, data []byte) error {
	// Create a local path: sites/host/path.
	localPath := filepath.Join("sites", u.Host, u.Path)

	// If a path ends with "/", add index.html.
	if strings.HasSuffix(u.Path, "/") || filepath.Ext(u.Path) == "" {
		localPath = filepath.Join(localPath, "index.html")
	}
	// Create folders if they don't already exist.
	dir := filepath.Dir(localPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directories %s: %w", localPath, err)
	} else {
		log.Printf("created %s\n", dir)
	}

	// Write data to a file.
	if err := os.WriteFile(localPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", localPath, err)
	} else {
		log.Printf("saved %s\n", filepath.Base(localPath))
	}

	return nil
}
