package cron

import (
	"testing"

	"github.com/kilgaloon/leprechaun/config"
)

var (
	iniFile  = "../tests/configs/config_regular.ini"
	path     = &iniFile
	cfgWrap  = config.NewConfigs()
	fakeCron = New("test", cfgWrap.New("test", *path))
)

func TestStop(t *testing.T) {
	fakeCron.Event.Subscribe("cron:ready", func() {
		fakeCron.Stop()
	})
}

func TestBuildJobs(t *testing.T) {
	fakeCron.buildJobs()
}

func TestStart(t *testing.T) {
	go fakeCron.Start()
}

func TestRegisterApiHandles(t *testing.T) {
	cmds := fakeCron.RegisterAPIHandles()
	if len(cmds) > 0 {
		t.Fail()
	}
}
