package client

import (
	"net/http"
	"net/http/httptest"
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
			if def.GetName() != "test" {
				t.Fatal("Agent name needs to be test")
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
				t.Fatal("Info command doesnt exist")
			}

			// lookup := 0
			// for {
			// 	if lookup > 50 {
			// 		t.Fatal("Lookup exceeded")
			// 		break
			// 	}

			// 	lookup++
			// 	if def.FindRecipe("test") == nil {
			// 		dat, err := ioutil.ReadFile(fakeClient.GetConfig().GetRecipesPathAbs() + "/../recipe.test")
			// 		if err != nil {
			// 			t.Fatal(err)
			// 		}

			// 		ioutil.WriteFile(fakeClient.GetConfig().GetRecipesPathAbs()+"/test.yml", dat, 0777)
			// 	} else {
			// 		os.Remove(fakeClient.GetConfig().GetRecipesPathAbs() + "/test.yml")

			// 		break
			// 	}

			// }

			if foo, ok := cmds["stop"]; ok {
				req, err := http.NewRequest("GET", "/scheduler/stop", nil)
				if err != nil {
					t.Fatal(err)
				}

				rr := httptest.NewRecorder()

				foo(rr, req)

				if rr.Code != http.StatusOK {
					t.Fatal("Expected code is 200")
				}

				if len(def.Queue.Stack) > 0 {
					t.Fatal("Agent queue stack still populated")
				}

				for {
					if fakeClient.GetStatus() == daemon.Stopped {
						rr := httptest.NewRecorder()
						foo(rr, req)
						if rr.Code == http.StatusOK {
							t.Fatal("Client is already stopped, status code 512 is expected")
						}

						break
					}
				}
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
