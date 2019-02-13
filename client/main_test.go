package client

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
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

func TestMain(t *testing.T) {
	wg.Add(3)
	go fakeClient.Start()
	go fakeClient2.Start()
}
func TestStop(t *testing.T) {
	fakeClient2.Event.Subscribe("client:ready", func() {
		dat, _ := ioutil.ReadFile(fakeClient2.GetConfig().GetRecipesPath() + "/schedule.yml")
		ioutil.WriteFile(fakeClient2.GetConfig().GetRecipesPath()+"/test.yml", dat, 0644)

		os.Remove(fakeClient2.GetConfig().GetRecipesPath() + "/test.yml")

		fakeClient2.GetMutex().Lock()
		defer fakeClient2.GetMutex().Unlock()

		if len(fakeClient2.Queue.Stack) < 1 {
			t.Error("Add to queue failed when client is running and new recipe is added")
		}

		_, err := fakeClient2.Stop(os.Stdin, "")
		if err != nil {
			t.Fail()
		}
	})
}
func TestLockUnlock(t *testing.T) {
	//event.EventHandler.Subscribe("client:ready", func() {
	fakeClient.Event.Subscribe("client:ready", func() {
		fakeClient.Lock()
		if !fakeClient.isWorking() {
			t.Fail()
			return
		}

		fakeClient.GetPID()

		fakeClient.Unlock()
	})
	//})
}

func TestRegisterAPIHandles(t *testing.T) {
	cmds := fakeClient.RegisterAPIHandles()

	fakeClient.Mu.Lock()
	defer fakeClient.Mu.Unlock()
	if foo, ok := cmds["info"]; ok {
		req, err := http.NewRequest("GET", "/client/info", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		foo(rr, req)
	} else {
		t.Fail()
	}

}
