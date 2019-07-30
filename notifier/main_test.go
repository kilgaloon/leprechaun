package notifier

import (
	"testing"

	notis "github.com/kilgaloon/leprechaun/notifier/notifications"

	"github.com/kilgaloon/leprechaun/config"
	"github.com/kilgaloon/leprechaun/log"
)

var (
	iniFile = "../tests/configs/config_regular.ini"
	path    = &iniFile
	cfgWrap = config.NewConfigs()
	cfg     = cfgWrap.New("test", *path)
	logger  = log.Logs{
		ErrorLog: cfg.ErrorLog(),
		InfoLog:  cfg.InfoLog(),
	}
	n = New(cfg, logger)
)

func TestNotifyWithOptions(t *testing.T) {
	n.Methods["mail"] = notis.NewEmail(cfg)

	n.NotifyWithOptions(notis.Options{
		Body: "test body",
	})
}
