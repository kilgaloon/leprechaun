package client

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/kilgaloon/leprechaun/config"
)

var (
	iniFile     = "../tests/configs/config_regular.ini"
	path        = &iniFile
	cfgWrap     = config.NewConfigs()
	fakeClient  = New("test", cfgWrap.New("test", *path), false)
	fakeClient2 = New("test", cfgWrap.New("test", *path), false)
)

func TestMain(t *testing.T) {
	go fakeClient.Start()
	go fakeClient2.Start()

	fakeClient2.Event.Subscribe("client:ready", func() {
		dat, _ := ioutil.ReadFile(fakeClient2.GetConfig().GetRecipesPath() + "/schedule.yml")
		ioutil.WriteFile(fakeClient2.GetConfig().GetRecipesPath()+"/test.yml", dat, 0644)

		os.Remove(fakeClient2.GetConfig().GetRecipesPath() + "/test.yml")

		fakeClient2.GetMutex().Lock()
		defer fakeClient2.GetMutex().Unlock()

		if len(fakeClient2.Queue.Stack) < 1 {
			t.Error("Add to queue failed when client is running and new recipe is added")
		}

		fakeClient2.Stop()
		if fakeClient2.stopped != true {
			t.Fail()
		}
	})

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

	cmds := fakeClient.RegisterAPIHandles()

	fakeClient.Mu.Lock()

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

	fakeClient.Mu.Unlock()
}
