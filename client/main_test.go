package client

import (
	"io/ioutil"
	"os"
	"strconv"
	"testing"

	"github.com/kilgaloon/leprechaun/config"
)

var (
	iniFile    = "../tests/configs/config_regular.ini"
	path       = &iniFile
	cfgWrap    = config.NewConfigs()
	fakeClient = New("test", cfgWrap.New("test", *path))
)

func TestStart(t *testing.T) {
	// remove hanging .lock file
	os.Remove(fakeClient.Agent.GetConfig().GetLockFile())
	// SetPID of client
	fakeClient.SetPID()
	// build queue

	// Fail because clien't isn't working anything on start
	if fakeClient.isWorking() {
		t.Error("Client should not be working anything here")
	}

	b, err := ioutil.ReadFile(fakeClient.Agent.GetConfig().GetPIDFile())
	if err != nil {
		t.Errorf("%s", err)
	}

	str := string(b)
	pid, err := strconv.Atoi(str)
	if err != nil {
		t.Error(err)
	}

	if fakeClient.GetPID() != pid {
		t.Errorf("PID expected to be %d but got %d", pid, fakeClient.GetPID())
	}
}
func TestLock(t *testing.T) {
	fakeClient.Lock()

	if !fakeClient.isWorking() {
		t.Fail()
	}

	fakeClient.Unlock()

	if fakeClient.isWorking() {
		t.Fail()
	}
}
