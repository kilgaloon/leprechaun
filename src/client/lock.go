package client

import (
	"os"
)

// LockProcess recipe on process
func LockProcess(name string, client *Client) (bool) {
	file := client.Config.recipesPath + "/" + name + ".lock"
	_, err := os.OpenFile(file, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return false
	}

	return true
}
// IsLocked check does lock exists
func IsLocked(name string, client *Client) (bool) {
	file := client.Config.recipesPath + "/" + name + ".lock"
	if _, err := os.Stat(file); err == nil {
		return true
	}

	return false
}
// RemoveLock recipe on process
func RemoveLock(name string, client *Client) (bool) {
	file := client.Config.recipesPath + "/" + name + ".lock"
	// delete file
	var err = os.Remove(file)
	if err != nil {
		return false
	}

	return true
}