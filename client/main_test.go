package client

import (
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
	t.Run("stop", func(t *testing.T) {
		fakeClient.Stop()
	})
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
