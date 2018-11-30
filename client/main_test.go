package client

import (
	"os"
	"sync"
	"testing"

	"github.com/kilgaloon/leprechaun/config"
)

var (
	iniFile     = "../tests/configs/config_regular.ini"
	path        = &iniFile
	cfgWrap     = config.NewConfigs()
	fakeClient  = New("test", cfgWrap.New("test", *path))
	fakeClient2 = New("test", cfgWrap.New("test", *path))
	wg          = new(sync.WaitGroup)
)

func TestStop(t *testing.T) {
	fakeClient2.Event.Subscribe("client:ready", func() {
		_, err := fakeClient2.Stop(os.Stdin, "")
		if err != nil {
			t.Fail()
		}
	})
	//})
}
func TestLockUnlock(t *testing.T) {
	//event.EventHandler.Subscribe("client:ready", func() {
	fakeClient.Event.Subscribe("client:ready", func() {
		fakeClient.Lock()
		if !fakeClient.isWorking() {
			t.Fail()
			return
		}

		fakeClient.Unlock()
	})
	//})
}
func TestStart(t *testing.T) {
	wg.Add(3)
	go fakeClient.Start()
	go fakeClient2.Start()
}
