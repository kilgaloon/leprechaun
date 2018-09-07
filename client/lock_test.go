package client

import (
	"os"
	"testing"
)

func TestLockProcess(t *testing.T) {
	LockProcess("recipe", fakeClient)

	file := fakeClient.Config.recipesPath + "/recipe.lock"
	if _, err := os.Stat(file); err != nil {
		t.Fail()
	}

	if !IsLocked("recipe", fakeClient) {
		t.Fail()
	}

	RemoveLock("recipe", fakeClient)

	if IsLocked("recipe", fakeClient) {
		t.Fail()
	}
}
