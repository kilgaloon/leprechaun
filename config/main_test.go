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
	assert.Equal(t, p, cfg.GetPath())
	assert.Equal(t, ErrorLog, cfg.GetErrorLog())
	assert.Equal(t, InfoLog, cfg.GetInfoLog())
	assert.Equal(t, RecipesPath, cfg.GetRecipesPath())
	assert.Equal(t, PIDFile, cfg.GetPIDFile())
	assert.Equal(t, LockFile, cfg.GetLockFile())
	assert.Equal(t, MaxAllowedWorkers, cfg.GetMaxAllowedWorkers())
	assert.Equal(t, MaxAllowedQueueWorkers, cfg.GetMaxAllowedQueueWorkers())
	assert.Equal(t, ServerPort, cfg.GetPort())
	assert.Equal(t, WorkerOutputDir, cfg.GetWorkerOutputDir())
	assert.Equal(t, "", cfg.GetNotificationsEmail())
	assert.Equal(t, "", cfg.GetSMTPHost())
	assert.Equal(t, "", cfg.GetSMTPUsername())
	assert.Equal(t, "", cfg.GetSMTPPassword())
	assert.Equal(t, []string([]string{"", "www."}), cfg.GetServerDomain())
}

func TestBuildGlobalFallback(t *testing.T) {
	cfg := ConfigGlobalFb

	p, _ := filepath.Abs("../tests/configs/config_global_fb.ini")
	assert.Equal(t, p, cfg.GetPath())
	assert.Equal(t, "../tests/var/log/leprechaun/error.log", cfg.GetErrorLog())
	assert.Equal(t, "../tests/var/log/leprechaun/info.log", cfg.GetInfoLog())
	assert.Equal(t, "../tests/etc/leprechaun/recipes", cfg.GetRecipesPath())
	assert.Equal(t, "../tests/var/run/leprechaun/.pid", cfg.GetPIDFile())
	assert.Equal(t, "../tests/var/run/leprechaun/.lock", cfg.GetLockFile())
	assert.Equal(t, "../tests/var/log/leprechaun/workers.output", cfg.GetWorkerOutputDir())
	assert.Equal(t, 5, cfg.GetMaxAllowedWorkers())
	assert.Equal(t, 5, cfg.GetMaxAllowedQueueWorkers())
	assert.Equal(t, "some@mail.com", cfg.GetNotificationsEmail())
	assert.Equal(t, "smtp.host.com", cfg.GetSMTPHost())
	assert.Equal(t, "smtp_user", cfg.GetSMTPUsername())
	assert.Equal(t, "smtp_pass", cfg.GetSMTPPassword())
}

func TestBuildWithSettings(t *testing.T) {
	cfg := ConfigWithSettings

	p, _ := filepath.Abs("../tests/configs/config_regular.ini")
	assert.Equal(t, p, cfg.GetPath())
	assert.Equal(t, "../tests/var/log/leprechaun/error.log", cfg.GetErrorLog())
	assert.Equal(t, "../tests/var/log/leprechaun/info.log", cfg.GetInfoLog())
	assert.Equal(t, "../tests/etc/leprechaun/recipes", cfg.GetRecipesPath())
	assert.Equal(t, "../tests/var/run/leprechaun/.pid", cfg.GetPIDFile())
	assert.Equal(t, "../tests/var/run/leprechaun/.lock", cfg.GetLockFile())
	assert.Equal(t, "../tests/var/log/leprechaun/workers.output", cfg.GetWorkerOutputDir())
	assert.Equal(t, 5, cfg.GetMaxAllowedWorkers())
	assert.Equal(t, 5, cfg.GetMaxAllowedQueueWorkers())
	assert.Equal(t, "some@mail.com", cfg.GetNotificationsEmail())
	assert.Equal(t, "smtp.host.com", cfg.GetSMTPHost())
	assert.Equal(t, "smtp_user", cfg.GetSMTPUsername())
	assert.Equal(t, "smtp_pass", cfg.GetSMTPPassword())
	assert.Equal(t, []string{"example.com", "www.example.com"}, cfg.GetServerDomain())
}

func TestBuildWithInvalidValues(t *testing.T) {
	cfg := ConfigWithInvalidValues

	assert.Equal(t, ErrorLog, cfg.GetErrorLog())
	assert.Equal(t, InfoLog, cfg.GetInfoLog())
	assert.Equal(t, RecipesPath, cfg.GetRecipesPath())
	assert.Equal(t, PIDFile, cfg.GetPIDFile())
	assert.Equal(t, LockFile, cfg.GetLockFile())
	assert.Equal(t, MaxAllowedWorkers, cfg.GetMaxAllowedWorkers())
	assert.Equal(t, MaxAllowedQueueWorkers, cfg.GetMaxAllowedQueueWorkers())
	assert.Equal(t, []string([]string{"", "www."}), cfg.GetServerDomain())
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
