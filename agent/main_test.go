package agent

import (
	"errors"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/kilgaloon/leprechaun/api"
	"github.com/kilgaloon/leprechaun/config"
)

var (
	iniFile      = "../tests/configs/config_regular.ini"
	path         = &iniFile
	cfgWrap      = config.NewConfigs()
	defaultAgent = New("test", cfgWrap.New("test", *path))
)

func TestGetterers(t *testing.T) {
	defaultAgent.GetName()
	defaultAgent.GetWorkers()
	defaultAgent.GetContext()
	defaultAgent.GetLogs()
	defaultAgent.GetConfig()
	defaultAgent.GetSocket()
	defaultAgent.GetMutex()
	defaultAgent.SetPID(1)
	if defaultAgent.GetPID() != 1 {
		t.Fail()
	}

	defaultAgent.GetStdout()
	defaultAgent.GetStdin()

	defaultAgent.SetStdin(os.Stdin)
	defaultAgent.SetStdout(os.Stdout)

	var a string
	fmt.Fprintf(defaultAgent, "%s", "Test write")
	fmt.Fscanf(defaultAgent, "%s", &a)

	externalCmds := make(map[string]api.Command)

	externalCmds["default:test"] = api.Command{
		Closure: func(r io.Writer, args ...string) ([][]string, error) {
			return [][]string{}, errors.New("Test error")
		},
		Definition: api.Definition{
			Text:  "Kills currently active worker by job name",
			Usage: "{agent} workers:kill {job}",
		},
	}

	defaultAgent.DefaultCommands(externalCmds)
}

func TestCommands(t *testing.T) {
	// no workers currently working
	defaultAgent.WorkersList(defaultAgent.GetStdout())

	// create worker
	_, err := defaultAgent.GetWorkers().CreateWorker("jobber")
	if err != nil {
		t.Fail()
	}

	defaultAgent.WorkersList(defaultAgent.GetStdout())
	// not existent worker
	defaultAgent.KillWorker(defaultAgent.GetStdout(), "test_job")
	defaultAgent.KillWorker(defaultAgent.GetStdout(), "jobber")
}
