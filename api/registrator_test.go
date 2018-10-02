package api

import (
	"testing"

	"github.com/kilgaloon/leprechaun/agent"

	"github.com/kilgaloon/leprechaun/config"
)

var (
	iniFile     = "../tests/configs/config_regular.ini"
	path        = &iniFile
	cfg         = config.NewConfigs()
	fakeClient  = agent.New("test", cfg.New("test", *path))
	registrator = CreateRegistrator(fakeClient)
)

func TestCommand(t *testing.T) {
	registrator.Command("test", func(args ...string) ([][]string, error) {
		resp := tabel{
			column{"test"},
		}

		return resp, nil
	})

	registrator.Command("test 2", func(args ...string) ([][]string, error) {
		resp := tabel{
			column{"test 2"},
		}

		return resp, nil
	})

	if len(registrator.Commands) < 2 {
		t.Fatalf("Failed to register all commands for registrator")
	}
}

func TestCall(t *testing.T) {
	registrator.Command("test", func(args ...string) ([][]string, error) {
		resp := tabel{
			column{"test"},
		}

		return resp, nil
	})

	resp, err := registrator.Call("test")
	if err != nil {
		t.Fail()
	}

	if resp[0][0] != "test" {
		t.Fail()
	}
}

func TestCallWithArguments(t *testing.T) {
	registrator.Command("test", func(args ...string) ([][]string, error) {
		resp := tabel{
			args,
		}

		return resp, nil
	})

	resp, err := registrator.Call("test", "arg1", "arg2")
	if err != nil {
		t.Fail()
	}

	if resp[0][0] != "arg1" || resp[0][1] != "arg2" {
		t.Fail()
	}
}

func TestCallCommandNotExist(t *testing.T) {
	_, err := registrator.Call("testability")
	if err == nil {
		t.Fail()
	}
}
