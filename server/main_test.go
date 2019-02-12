package server

import (
	"errors"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/kilgaloon/leprechaun/config"
)

var (
	iniFile     = "../tests/configs/config_regular.ini"
	path        = &iniFile
	cfgWrap     = config.NewConfigs()
	fakeServer  = New("test", cfgWrap.New("test", *path))
	fakeServer2 = New("test", cfgWrap.New("test", *path))
)

func TestStartStop(t *testing.T) {
	go fakeServer.Start()
	// retry 5 times before failing
	// this means server failed to start
	port := strconv.Itoa(fakeServer.GetConfig().GetPort())
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

		fakeServer.Stop()
		break
	}
}

func TestFindInPool(t *testing.T) {
	fakeServer.BuildPool()
	fakeServer.FindInPool("223344")

	recipe := fakeServer.Pool.Stack["223344"]
	recipe.Err = errors.New("Some random error")

	fakeServer.FindInPool("223344")

	fakeServer.BuildPool()
}

func TestIsTLS(t *testing.T) {
	if fakeServer.isTLS() {
		t.Fail()
	}
}
