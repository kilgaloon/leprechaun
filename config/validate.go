package config

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

// IsFileValid checks is path really file
func IsFileValid(path string, ext string) bool {
	if filepath.Ext(path) != ext {
		return false
	}
	// Check if file already exists
	if _, err := os.Stat(path); err == nil {
		return true
	}

	// Attempt to create it
	var d []byte
	if err := ioutil.WriteFile(path, d, 0644); err == nil {
		os.Remove(path) // And delete it
		return true
	}

	return false
}

// IsDirValid checks is path really file
func IsDirValid(path string) bool {
	// Get file stat
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	// check is dir
	if info.IsDir() {
		return true
	}

	return false
}
