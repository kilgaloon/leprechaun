package client

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/kilgaloon/leprechaun/event"

	"github.com/kilgaloon/leprechaun/config"
)

var (
	iniFile    = "../tests/configs/config_regular.ini"
	path       = &iniFile
	cfgWrap    = config.NewConfigs()
	fakeClient = New("test", cfgWrap.New("test", *path))
)

func TestStart(t *testing.T) {
	go fakeClient.Start()
}

func TestStop(t *testing.T) {
	tmpfile, err := ioutil.TempFile("./", "agent.stdin")
	if err != nil {
		log.Fatal(err)
	}

	defer os.Remove(tmpfile.Name()) // clean up

	if _, err := tmpfile.Write([]byte("Y\n")); err != nil {
		log.Fatal(err)
	}

	if _, err := tmpfile.Seek(0, 0); err != nil {
		log.Fatal(err)
	}

	fakeClient.Agent.SetStdin(tmpfile)

	fakeClient.Lock()
	fakeClient.Stop()

	tmpfile.Close()
}
func TestLockUnlock(t *testing.T) {
	fakeClient.Lock()
	if !fakeClient.isWorking() {
		t.Fail()
	}
	event.EventHandler.Dispatch("client:unlock")
}

func TestGetAgent(t *testing.T) {
	if fakeClient.GetAgent() != fakeClient.Agent {
		t.Fail()
	}
}
