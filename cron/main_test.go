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

func TestStart(t *testing.T) {
	go fakeCron.Start()

	fakeCron.buildJobs()
}

func TestRegisterCommands(t *testing.T) {
	fakeCron.RegisterCommands()
}

func TestStop(t *testing.T) {
	fakeCron.Stop()
}
