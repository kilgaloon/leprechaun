package notifications

import (
	"testing"

	"github.com/kilgaloon/leprechaun/config"
)

var (
	iniFile = "../../tests/configs/config_regular.ini"
	path    = &iniFile
	cfgWrap = config.NewConfigs()
	cfg     = cfgWrap.New("test", *path)
	n       = NewEmail(cfg)
)

func TestSettings(t *testing.T) {
	n.SetMessage("test")
	if n.GetMessage() != "test" {
		t.Fail()
	}

	n.SetTitle("test")
	if n.GetTitle() != "test" {
		t.Fail()
	}

	n.Send()
}
