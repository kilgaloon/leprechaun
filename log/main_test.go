package log

import (
	"os"
	"testing"

	"github.com/kilgaloon/leprechaun/config"
)

var (
	cfg    = config.BuildConfig("../tests/configs/config_regular.ini")
	logger = Logs{
		ErrorLog: cfg.ClientConfig.ErrorLog,
		InfoLog:  cfg.ClientConfig.InfoLog,
	}
)

func TestErrorLog(t *testing.T) {
	// log some random error message
	logger.Error("Some error message")

	info, err := os.Stat(logger.ErrorLog)
	if err != nil {
		t.Errorf("Failed because %s", err)
	}

	if !(info.Size() > 0) {
		t.Errorf("Filesize expected to be larger the 0, got %d", info.Size())
	}
	// first remove file
	os.Remove(logger.ErrorLog)

}

func TestInfoLog(t *testing.T) {
	// log some random error message
	logger.Info("Some info message")

	info, err := os.Stat(logger.InfoLog)
	if err != nil {
		t.Errorf("Failed because %s", err)
	}

	if !(info.Size() > 0) {
		t.Errorf("Filesize expected to be larger the 0, got %d", info.Size())
	}
	// first remove file
	os.Remove(logger.InfoLog)

}
