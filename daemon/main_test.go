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

func TestMain(t *testing.T) {
	// init daemon
	Init()
}

func TestAddService(t *testing.T) {
	d := Init()

	d.Cmd = "run fake_service"
	d.AddService(&fakeService{})

	go d.Run()
}

func TestRun(t *testing.T) {
	d := Init()
	go d.Run()

	if d.GetPID() != os.Getpid() {
		t.Fatal("PID NOT MATCHED")
	}

	req, err := http.NewRequest("GET", "/daemon/info", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	d.daemonInfo(rr, req)

	for i := 0; i < 5; i++ {
		_, err := http.Get("http://localhost:11401")
		if err != nil {
			// handle error
			time.Sleep(2 * time.Second)
			continue
		}

		d.GetInfo()
		d.renderInfo()

	}
}

func TestRunningDaemonInfo(t *testing.T) {
	d := Init()

	d.Cmd = "daemon info"
	d.Run()
	for i := 0; i < 5; i++ {
		_, err := http.Get("http://localhost:11401")
		if err != nil {
			// handle error
			d.API.Start()
			time.Sleep(2 * time.Second)
			continue
		}

		break
	}
}

func TestRunningDaemonServices(t *testing.T) {
	d := Init()

	d.Cmd = "daemon services"
	d.Run()
	for i := 0; i < 5; i++ {
		_, err := http.Get("http://localhost:11401")
		if err != nil {
			// handle error
			d.API.Start()
			time.Sleep(2 * time.Second)
			continue
		}

		break
	}
}

func TestRunningDaemonKill(t *testing.T) {
	d := Init()

	d.Cmd = "daemon kill"
	d.Run()
	for i := 0; i < 5; i++ {
		_, err := http.Get("http://localhost:11401")
		if err != nil {
			// handle error
			d.API.Start()
			time.Sleep(2 * time.Second)
			continue
		}

		break
	}
}

func TestAPIRunning(t *testing.T) {
	os.Setenv("RUN_MODE", "test")

	d := Init()

	d.Cmd = ""
	d.Run()

	d = Init()

	d.Kill()
}

func TestHelpFlag(t *testing.T) {
	helpCommands()
}
