package cron

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kilgaloon/leprechaun/config"
	"github.com/kilgaloon/leprechaun/daemon"
	"github.com/kilgaloon/leprechaun/recipe"
)

var (
	iniFile  = "../tests/configs/config_regular.ini"
	path     = &iniFile
	cfgWrap  = config.NewConfigs()
	def      = &Cron{}
	fakeCron = def.New("test", cfgWrap.New("test", *path), false)
)

func TestMain(t *testing.T) {
	if fakeCron.GetName() != "test" {
		t.Fail()
	}
}
func TestBuildJobs(t *testing.T) {
	Agent.buildJobs()
}

func TestStartStop(t *testing.T) {
	go fakeCron.Start()

	for {
		if fakeCron.GetStatus() == daemon.Started {

			cmds := fakeCron.RegisterAPIHandles()

			if foo, ok := cmds["info"]; ok {
				req, err := http.NewRequest("GET", "/cron/info", nil)
				if err != nil {
					t.Fatal(err)
				}

				rr := httptest.NewRecorder()

				foo(rr, req)
			} else {
				t.Fail()
			}

			if foo, ok := cmds["stop"]; ok {
				req, err := http.NewRequest("GET", "/cron/stop", nil)
				if err != nil {
					t.Fatal(err)
				}

				rr := httptest.NewRecorder()

				foo(rr, req)
			} else {
				t.Fail()
			}

			if foo, ok := cmds["pause"]; ok {
				req, err := http.NewRequest("GET", "/cron/pause", nil)
				if err != nil {
					t.Fatal(err)
				}

				rr := httptest.NewRecorder()

				foo(rr, req)
			} else {
				t.Fail()
			}

			if foo, ok := cmds["start"]; ok {
				req, err := http.NewRequest("GET", "/cron/start", nil)
				if err != nil {
					t.Fatal(err)
				}

				rr := httptest.NewRecorder()

				foo(rr, req)
			} else {
				t.Fail()
			}

			// test that cron is started again if its already started
			if foo, ok := cmds["start"]; ok {
				req, err := http.NewRequest("GET", "/cron/start", nil)
				if err != nil {
					t.Fatal(err)
				}

				rr := httptest.NewRecorder()

				foo(rr, req)
			}

			fakeCron.Stop()
			if fakeCron.GetStatus() != daemon.Stopped {
				t.Fail()
			}

			break
		}
	}
}

func TestRegisterApiHandles(t *testing.T) {
	cmds := fakeCron.RegisterAPIHandles()
	if len(cmds) > 4 {
		t.Fail()
	}
}

func TestPrepareAndRun(t *testing.T) {
	r, err := recipe.Build("../tests/etc/leprechaun/recipes/cron.yml")
	if err != nil {
		t.Fail()
	}

	Agent.prepareAndRun(r)
}
