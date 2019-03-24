package daemon

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/kilgaloon/leprechaun/config"
)

var (
	configs            = config.NewConfigs()
	ConfigWithSettings = configs.New("test", "../tests/configs/config_regular.ini")
)

type fakeService struct {
}

func (fs fakeService) GetName() string {
	return "fake_service"
}

// RegisterAPIHandles to be used in socket communication
// If you want to takeover default commands from agent
// call DefaultCommands from Agent which is same command
func (fs *fakeService) RegisterAPIHandles() map[string]func(w http.ResponseWriter, r *http.Request) {
	cmds := make(map[string]func(w http.ResponseWriter, r *http.Request))

	return cmds
}

func (fs *fakeService) DefaultAPIHandles() map[string]func(w http.ResponseWriter, r *http.Request) {
	cmds := make(map[string]func(w http.ResponseWriter, r *http.Request))

	return cmds
}

func (fs *fakeService) GetStatus() ServiceStatus {
	return ServiceStatus(1)
}

func (fs *fakeService) GetConfig() config.AgentConfig {
	return *ConfigWithSettings
}

func (fs *fakeService) IsDebug() bool {
	return true
}

func (fs *fakeService) SetStatus(s int)         {}
func (fs *fakeService) Start()                  {}
func (fs *fakeService) Stop()                   {}
func (fs *fakeService) Pause()                  {}
func (fs *fakeService) SetPipeline(chan string) {}
func (fs *fakeService) New(name string, cfg *config.AgentConfig, debug bool) Service {
	srv := &fakeService{}
	return srv
}

func TestAddService(t *testing.T) {
	Srv.AddService(&fakeService{})
}

func TestRun(t *testing.T) {
	go Srv.Run()

	if Srv.GetPID() != os.Getpid() {
		t.Fatal("PID NOT MATCHED")
	}

	req, err := http.NewRequest("GET", "/daemon/info", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	Srv.daemonInfo(rr, req)

	for i := 0; i < 5; i++ {
		_, err := http.Get("http://localhost:11401")
		if err != nil {
			// handle error
			time.Sleep(2 * time.Second)
			continue
		}

		Srv.GetInfo()
		Srv.renderInfo()
		Srv.killDaemon()
		Srv.daemonServices()

	}

	//Srv.Kill()
}

func TestRunningDaemonInfo(t *testing.T) {
	Srv.Cmd = "daemon info"
	Srv.Run()
	for i := 0; i < 5; i++ {
		_, err := http.Get("http://localhost:11401")
		if err != nil {
			// handle error
			Srv.API.Start()
			time.Sleep(2 * time.Second)
			continue
		}

		break
	}
}

func TestRunningDaemonServices(t *testing.T) {
	Srv.Cmd = "daemon services"
	Srv.Run()
	for i := 0; i < 5; i++ {
		_, err := http.Get("http://localhost:11401")
		if err != nil {
			// handle error
			Srv.API.Start()
			time.Sleep(2 * time.Second)
			continue
		}

		break
	}
}

func TestRunningDaemonKill(t *testing.T) {
	Srv.Cmd = "daemon kill"
	Srv.Run()
	for i := 0; i < 5; i++ {
		_, err := http.Get("http://localhost:11401")
		if err != nil {
			// handle error
			Srv.API.Start()
			time.Sleep(2 * time.Second)
			continue
		}

		break
	}
}
