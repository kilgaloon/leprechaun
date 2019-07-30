package server

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/kilgaloon/leprechaun/config"
)

var (
	iniFile     = "../tests/configs/config_regular.ini"
	path        = &iniFile
	cfgWrap     = config.NewConfigs()
	def         = &Server{}
	fakeServer  = def.New("test", cfgWrap.New("test", *path), false)
	fakeServer2 = def.New("test", cfgWrap.New("test", *path), false)
)

func TestStartStop(t *testing.T) {
	if fakeServer.GetName() != "test" {
		t.Fail()
	}

	go fakeServer.Start()
	// retry 5 times before failing
	// this means server failed to start
	port := strconv.Itoa(fakeServer.GetConfig().Port())
	for i := 0; i < 5; i++ {
		_, err := http.Get("http://localhost" + ":" + port)
		if err != nil {
			// handle error
			time.Sleep(2 * time.Second)
			continue
		}

		TestFindInPool(t)

		_, err = http.Get("http://localhost" + ":" + port + "/ping")
		if err != nil {
			t.Fail()
		}

		_, err = http.Get("http://localhost" + ":" + port + "/hook?id=223344")
		if err != nil {
			t.Fail()
		}

		cmds := fakeServer.RegisterAPIHandles()

		if foo, ok := cmds["info"]; ok {
			req, err := http.NewRequest("GET", "/server/info", nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()

			foo(rr, req)
		} else {
			t.Fail()
		}

		if foo, ok := cmds["stop"]; ok {
			req, err := http.NewRequest("GET", "/server/stop", nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()

			foo(rr, req)
		} else {
			t.Fail()
		}

		if foo, ok := cmds["pause"]; ok {
			req, err := http.NewRequest("GET", "/server/pause", nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()

			foo(rr, req)
		} else {
			t.Fail()
		}

		if foo, ok := cmds["start"]; ok {
			req, err := http.NewRequest("GET", "/server/start", nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()

			foo(rr, req)
		} else {
			t.Fail()
		}

		fakeServer.Stop()
		break
	}

	// Agent.Lock()
	// fakeServer2.GetConfig().Domain = "https://localhost"
	// Agent.Unlock()

	go fakeServer2.Start()

}

func TestFindInPool(t *testing.T) {
	Agent.BuildPool()
	Agent.FindInPool("223344")

	Agent.Lock()
	recipe := Agent.Pool.Stack["223344"]
	recipe.Err = errors.New("Some random error")
	Agent.Unlock()

	Agent.FindInPool("223344")

	Agent.BuildPool()
}

func TestIsTLS(t *testing.T) {
	// Agent.Lock()
	// fakeServer2.GetConfig().Domain = "localhost"
	// Agent.Unlock()

	if Agent.isTLS() {
		t.Fail()
	}
}

func TestRegisterAPIHandles(t *testing.T) {
	cmds := fakeServer.RegisterAPIHandles()
	if len(cmds) > 4 {
		t.Fail()
	}
}
