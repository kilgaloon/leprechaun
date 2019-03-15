package client

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/kilgaloon/leprechaun/config"
	"github.com/kilgaloon/leprechaun/daemon"
)

var (
	iniFile    = "../tests/configs/config_regular.ini"
	path       = &iniFile
	cfgWrap    = config.NewConfigs()
	def        = &Client{}
	fakeClient = def.New("test", cfgWrap.New("test", *path), false)
)

func TestMain(t *testing.T) {
	go fakeClient.Start()

	for {
		if fakeClient.GetStatus() == daemon.Started {
			//lookup := 0
			for {
				// if lookup > 50 {
				// 	t.Fatal("Lookup exceeded")
				// 	break
				// }

				// lookup++
				if Agent.FindRecipe("test") == nil {
					dat, _ := ioutil.ReadFile(fakeClient.GetConfig().GetRecipesPathAbs() + "/../recipe.test")
					ioutil.WriteFile(fakeClient.GetConfig().GetRecipesPathAbs()+"/test.yml", dat, 0777)
				} else {
					os.Remove(fakeClient.GetConfig().GetRecipesPathAbs() + "/test.yml")

					break
				}

			}

			if Agent.GetName() != "test" {
				t.Fail()
			}

			cmds := fakeClient.RegisterAPIHandles()

			if foo, ok := cmds["info"]; ok {
				req, err := http.NewRequest("GET", "/scheduler/info", nil)
				if err != nil {
					t.Fatal(err)
				}

				rr := httptest.NewRecorder()

				foo(rr, req)
			} else {
				t.Fail()
			}

			fakeClient.Stop()
			if fakeClient.GetStatus() != daemon.Stopped {
				t.Fail()
			}

			if len(Agent.Queue.Stack) > 0 {
				t.Fail()
			}

			if foo, ok := cmds["stop"]; ok {
				req, err := http.NewRequest("GET", "/scheduler/stop", nil)
				if err != nil {
					t.Fatal(err)
				}

				rr := httptest.NewRecorder()

				foo(rr, req)
			} else {
				t.Fail()
			}

			if foo, ok := cmds["pause"]; ok {
				req, err := http.NewRequest("GET", "/scheduler/pause", nil)
				if err != nil {
					t.Fatal(err)
				}

				rr := httptest.NewRecorder()

				foo(rr, req)
			} else {
				t.Fail()
			}

			if foo, ok := cmds["start"]; ok {
				req, err := http.NewRequest("GET", "/scheduler/start", nil)
				if err != nil {
					t.Fatal(err)
				}

				rr := httptest.NewRecorder()

				foo(rr, req)
			} else {
				t.Fail()
			}

			break
		}
	}
}
