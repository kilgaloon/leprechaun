package client

import (
	"io/ioutil"
	"strconv"
	"testing"
)

var (
	iniFile    = "../tests/configs/client.ini"
	path       = &iniFile
	fakeClient = CreateAgent(path)
)

func TestStart(t *testing.T) {
	// remove hanging .lock file
	fakeClient.Unlock()
	// SetPID of client
	fakeClient.SetPID()
	// build queue

	// Fail because clien't isn't working anything on start
	if fakeClient.isWorking() {
		t.Fail()
	}

	b, err := ioutil.ReadFile(fakeClient.Config.PIDFile)
	if err != nil {
		t.Errorf("%s", err)
	}

	str := string(b)
	pid, err := strconv.Atoi(str)
	if err != nil {
		t.Fail()
	}
	
	if fakeClient.GetPID() != pid {
		t.Fail()
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
