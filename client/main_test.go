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

func TestStart(t *testing.T) {
	wg.Add(3)
	go fakeClient.Start()
	go fakeClient2.Start()
}

func TestStop(t *testing.T) {
	//event.EventHandler.Subscribe("client:ready", func() {
	go func() {
		for {
			select {
			case <-fakeClient2.ReadyChan:
				_, err := fakeClient2.Stop(os.Stdin, "")
				if err != nil {
					t.Fail()
				}

				break
			}
		}
	}()

	//})
}
func TestLockUnlock(t *testing.T) {
	//event.EventHandler.Subscribe("client:ready", func() {
	go func() {
		for {
			select {
			case <-fakeClient.ReadyChan:
				fakeClient.Lock()
				if !fakeClient.isWorking() {
					t.Fail()
				}

				fakeClient.Unlock()
				break
			}
		}
	}()

	//})
}
