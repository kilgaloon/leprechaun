package server

import (
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/kilgaloon/leprechaun/config"
)

var (
	iniFile    = "../tests/configs/config_regular.ini"
	path       = &iniFile
	cfgWrap    = config.NewConfigs()
	fakeServer = New("test", cfgWrap.New("test", *path))
)

func TestStartStop(t *testing.T) {
	go fakeServer.Start()
	// retry 5 times before failing
	// this means server failed to start
	port := strconv.Itoa(fakeServer.Agent.GetConfig().GetPort())
	for i := 0; i < 5; i++ {
		_, err := http.Get("http://localhost" + ":" + port)
		if err != nil {
			// handle error
			time.Sleep(2 * time.Second)
			continue
		}

		TestFindInPool(t)
		TestProcessRecipe(t)

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

func TestRegisterCommands(t *testing.T) {
	fakeServer.RegisterCommands()
}

func TestFindInPool(t *testing.T) {
	// simulate exceeding maximum number of workers
	for i := 0; i < fakeServer.Agent.GetConfig().GetMaxAllowedWorkers(); i++ {
		fakeServer.FindInPool("223344")
	}
}

func TestProcessRecipe(t *testing.T) {
	recipe := fakeServer.Pool.Stack["223344"]
	fakeServer.ProcessRecipe(recipe)
}
