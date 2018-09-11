package config

import (
	"io/ioutil"
	"os"
)

// IsFileValid checks is path really file
func IsFileValid(path string) bool {
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
