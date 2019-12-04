package config

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	configs                      = NewConfigs()
	ConfigWithoutDefaultSettings = configs.New("test", "../tests/configs/config_without_default_values.ini")
	ConfigWithSettings           = configs.New("test", "../tests/configs/config_regular.ini")
	ConfigWithInvalidValues      = configs.New("test", "../tests/configs/config_wrong_value.ini")
	ConfigWithWrongExt           = configs.New("test", "../tests/configs/config_wrong_ext.ini")
	ConfigGlobalFb               = configs.New("test", "../tests/configs/config_global_fb.ini")
)

func TestBuildWithoutSettings(t *testing.T) {
	cfg := ConfigWithoutDefaultSettings

	p, _ := filepath.Abs("../tests/configs/config_without_default_values.ini")
	rp, _ := filepath.Abs(RecipesPath)
	assert.Equal(t, p, cfg.GetPath())
	assert.Equal(t, ErrorLog, cfg.GetErrorLog())
	assert.Equal(t, InfoLog, cfg.GetInfoLog())
	assert.Equal(t, RecipesPath, cfg.GetRecipesPath())
	assert.Equal(t, rp, cfg.GetRecipesPathAbs())
	assert.Equal(t, MaxAllowedWorkers, cfg.GetMaxAllowedWorkers())
	assert.Equal(t, MaxAllowedQueueWorkers, cfg.GetMaxAllowedQueueWorkers())
	assert.Equal(t, ServerPort, cfg.GetPort())
	assert.Equal(t, WorkerOutputDir, cfg.GetWorkerOutputDir())
	assert.Equal(t, "", cfg.GetNotificationsEmail())
	assert.Equal(t, "", cfg.GetSMTPHost())
	assert.Equal(t, "", cfg.GetSMTPUsername())
	assert.Equal(t, "", cfg.GetSMTPPassword())
	assert.Equal(t, []string([]string{"localhost", "www.localhost"}), cfg.GetServerDomain())
	assert.Equal(t, ErrorReporting, cfg.GetErrorReporting())
	assert.Equal(t, "", cfg.GetCertPemPath())
	assert.Equal(t, "", cfg.GetCertKeyPath())
	assert.Equal(t, make(map[string]string), cfg.GetRemoteServices())
	assert.Equal(t, "native", cfg.GetShell())
}

func TestBuildGlobalFallback(t *testing.T) {
	cfg := ConfigGlobalFb

	p, _ := filepath.Abs("../tests/configs/config_global_fb.ini")
	rp, _ := filepath.Abs("../tests/etc/leprechaun/recipes")
	assert.Equal(t, p, cfg.GetPath())
	assert.Equal(t, "../tests/var/log/leprechaun/error.log", cfg.GetErrorLog())
	assert.Equal(t, "../tests/var/log/leprechaun/info.log", cfg.GetInfoLog())
	assert.Equal(t, "../tests/etc/leprechaun/recipes", cfg.GetRecipesPath())
	assert.Equal(t, rp, cfg.GetRecipesPathAbs())
	assert.Equal(t, "../tests/var/log/leprechaun/workers.output", cfg.GetWorkerOutputDir())
	assert.Equal(t, 5, cfg.GetMaxAllowedWorkers())
	assert.Equal(t, 5, cfg.GetMaxAllowedQueueWorkers())
	assert.Equal(t, "some@mail.com", cfg.GetNotificationsEmail())
	assert.Equal(t, "smtp.host.com", cfg.GetSMTPHost())
	assert.Equal(t, "smtp_user", cfg.GetSMTPUsername())
	assert.Equal(t, "smtp_pass", cfg.GetSMTPPassword())
	assert.Equal(t, false, cfg.GetErrorReporting())
	assert.Equal(t, "cert.pem", cfg.GetCertPemPath())
	assert.Equal(t, "key.pem", cfg.GetCertKeyPath())

	rs := make(map[string]string)
	rs["localhost"] = "11000"
	rs["digioc"] = "11001"

	assert.Equal(t, rs, cfg.GetRemoteServices())
	assert.Equal(t, "bash", cfg.GetShell())
	
}

func TestBuildWithSettings(t *testing.T) {
	cfg := ConfigWithSettings

	p, _ := filepath.Abs("../tests/configs/config_regular.ini")
	rp, _ := filepath.Abs("../tests/etc/leprechaun/recipes")
	assert.Equal(t, p, cfg.GetPath())
	assert.Equal(t, "../tests/var/log/leprechaun/error.log", cfg.GetErrorLog())
	assert.Equal(t, "../tests/var/log/leprechaun/info.log", cfg.GetInfoLog())
	assert.Equal(t, "../tests/etc/leprechaun/recipes", cfg.GetRecipesPath())
	assert.Equal(t, rp, cfg.GetRecipesPathAbs())
	assert.Equal(t, "../tests/var/log/leprechaun/workers.output", cfg.GetWorkerOutputDir())
	assert.Equal(t, 5, cfg.GetMaxAllowedWorkers())
	assert.Equal(t, 5, cfg.GetMaxAllowedQueueWorkers())
	assert.Equal(t, "some@mail.com", cfg.GetNotificationsEmail())
	assert.Equal(t, "smtp.host.com", cfg.GetSMTPHost())
	assert.Equal(t, "smtp_user", cfg.GetSMTPUsername())
	assert.Equal(t, "smtp_pass", cfg.GetSMTPPassword())
	assert.Equal(t, []string{"example.com", "www.example.com"}, cfg.GetServerDomain())
	assert.Equal(t, true, cfg.GetErrorReporting())
	assert.Equal(t, "../tests/crts/certificate.pem", cfg.GetCertPemPath())
	assert.Equal(t, "../tests/crts/key.pem", cfg.GetCertKeyPath())

	rs := make(map[string]string)
	rs["localhost"] = "11400"

	assert.Equal(t, rs, cfg.GetRemoteServices())
	assert.Equal(t, "bash", cfg.GetShell())
}

func TestBuildWithInvalidValues(t *testing.T) {
	cfg := ConfigWithInvalidValues

	assert.Equal(t, ErrorLog, cfg.GetErrorLog())
	assert.Equal(t, InfoLog, cfg.GetInfoLog())
	assert.Equal(t, RecipesPath, cfg.GetRecipesPath())
	assert.Equal(t, MaxAllowedWorkers, cfg.GetMaxAllowedWorkers())
	assert.Equal(t, MaxAllowedQueueWorkers, cfg.GetMaxAllowedQueueWorkers())
	assert.Equal(t, []string([]string{"localhost", "www.localhost"}), cfg.GetServerDomain())
}

func TestGettingNotExistingConfig(t *testing.T) {
	// this should not break, instead we need to get empty object
	configs.GetConfig("not_exists")
}

func TestGettingExistingConfig(t *testing.T) {
	configs.GetConfig("test")
}

func TestNotValidPathToConfig(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			//
		}
	}()

	NewConfigs().New("test", "some_path")
}
