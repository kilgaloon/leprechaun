package log

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/kilgaloon/leprechaun/config"
)

var (
	cfgWrap = config.NewConfigs()
	cfg     = cfgWrap.New("test", "../tests/configs/config_regular.ini")
	cfg2    = cfgWrap.New("test", "../tests/configs/config_wrong_value.ini")
	logger  = Logs{
		Debug:    true,
		ErrorLog: cfg.ErrorLog(),
		InfoLog:  cfg.InfoLog(),
	}

	logger2 = Logs{
		Debug:    true,
		ErrorLog: cfg2.ErrorLog(),
		InfoLog:  cfg2.InfoLog(),
	}
)

func TestErrorLog(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()
	// log some random error message
	logger.Error("Some error message")
	logger2.Error("Some error message")

	info, err := os.Stat(logger.ErrorLog)
	if err != nil {
		t.Errorf("Failed because %s", err)
	}

	if info.Size() <= 0 {
		t.Errorf("Filesize expected to be larger the 0, got %d", info.Size())
	}
	// first remove file
	os.Remove(logger.ErrorLog)
	var d []byte
	ioutil.WriteFile(logger.ErrorLog, d, 0644)

}

func TestInfoLog(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()
	// log some random error message
	logger.Info("Some info message")
	logger2.Info("Some info message")

	info, err := os.Stat(logger.InfoLog)
	if err != nil {
		t.Errorf("Failed because %s", err)
	}

	if info.Size() <= 0 {
		t.Errorf("Filesize expected to be larger the 0, got %d", info.Size())
	}
	// first remove file
	os.Remove(logger.InfoLog)
	var d []byte
	ioutil.WriteFile(logger.InfoLog, d, 0644)
}
