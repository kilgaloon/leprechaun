package remote

import (
	"crypto/rand"
	"crypto/tls"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/kilgaloon/leprechaun/config"
	"github.com/kilgaloon/leprechaun/context"
	"github.com/kilgaloon/leprechaun/daemon"
	"github.com/kilgaloon/leprechaun/log"
	"github.com/kilgaloon/leprechaun/recipe"
	"github.com/kilgaloon/leprechaun/workers"
)

var (
	iniFile        = "../tests/configs/config_regular.ini"
	iniFileDebug   = "../tests/configs/config_regular_not_debug.ini"
	path           = &iniFile
	pathNotDebug   = &iniFileDebug
	cfgWrap        = config.NewConfigs()
	def            = &Remote{}
	debugDef       = &Remote{}
	fakeClient     = def.New("test", cfgWrap.New("test", *path), true)
	clientNotDebug = debugDef.New("test", cfgWrap.New("test", *pathNotDebug), false)
	r, err         = recipe.Build("../tests/etc/leprechaun/recipes/remote.yml")
	wrks           = workers.New(
		cfgWrap.New("test", *path),
		log.Logs{},
		context.New(),
		true,
	)
	worker, errr = wrks.CreateWorker(r)
)

func TestMain(t *testing.T) {
	go fakeClient.Start()

	for {
		if fakeClient.GetStatus() == daemon.Started {
			cmds := fakeClient.RegisterAPIHandles()

			if foo, ok := cmds["stop"]; ok {
				req, err := http.NewRequest("GET", "/remote/stop", nil)
				if err != nil {
					t.Fatal(err)
				}

				rr := httptest.NewRecorder()

				foo(rr, req)

				if rr.Code != http.StatusOK {
					t.Fatal("Expected code is 200")
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

			if foo, ok := cmds["start"]; ok {
				req, err := http.NewRequest("GET", "/remote/start", nil)
				if err != nil {
					t.Fatal(err)
				}

				rr := httptest.NewRecorder()

				foo(rr, req)

				if rr.Code != http.StatusOK {
					t.Fatal("Expected code is 200")
				}

				for {
					if fakeClient.GetStatus() == daemon.Started {
						rr := httptest.NewRecorder()
						foo(rr, req)
						if rr.Code == http.StatusOK {
							t.Fatal("Client is already started, status code 512 is expected")
						}

						break
					}
				}
			} else {
				t.Fatal("Start command not registered")
			}

			go worker.Run()
			break
		}
	}

	// clean
	for {
		_, err := os.Stat("remote.txt")
		if os.IsNotExist(err) {
			continue
		}

		cont, err := ioutil.ReadFile("remote.txt")
		if err != nil {
			continue
		}

		msg := string(cont)
		if msg == "" {
			continue
		}

		if msg != "testing this recipe\n" {
			t.Fatalf("Content of file doesn't match: %s", msg)
		}

		os.Remove("remote.txt")
		break
	}

	// clean
	for {
		_, err := os.Stat("buffer_test.txt")
		if os.IsNotExist(err) {
			continue
		}

		cont, err := ioutil.ReadFile("buffer_test.txt")
		if err != nil {
			continue
		}

		contReal, err := ioutil.ReadFile("../tests/etc/leprechaun/buffer.txt")
		if err != nil {
			continue
		}

		msg := string(cont)
		msgReal := string(contReal)
		if msg == "" || msgReal == "" {
			continue
		}

		if msg != msgReal {
			t.Fatalf("Content of file doesn't match: %s", msg)
		}

		os.Remove("buffer_test.txt")
		break
	}

}

func TestNotDebug(t *testing.T) {
	go clientNotDebug.Start()

	cert, err := tls.LoadX509KeyPair(
		clientNotDebug.GetConfig().GetCertPemPath(),
		clientNotDebug.GetConfig().GetCertKeyPath(),
	)

	if err != nil {
		t.Fatal(err)
	}

	config := tls.Config{Certificates: []tls.Certificate{cert}, InsecureSkipVerify: true}
	config.Rand = rand.Reader

	port := strconv.Itoa(clientNotDebug.GetConfig().GetPort())
	host := net.JoinHostPort("localhost", port)

	lookup := 0
	for {
		if lookup == 5 {
			t.Fatalf("Lookup failed")
		}

		conn, err := tls.Dial("tcp", host, &config)
		if err != nil {
			lookup++

			time.Sleep(5 * time.Second)

			continue
		}

		message := make([]byte, 5)
		// listen for message
		_, err = conn.Read(message)
		if err != nil {
			t.Fatal(err)
		}
		msg := string(message)
		if msg == "ready" {
			break
		}

		t.Fail()
	}
}
