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
	assert.Equal(t, p, cfg.Path())
	assert.Equal(t, ErrorLog, cfg.ErrorLog())
	assert.Equal(t, InfoLog, cfg.InfoLog())
	assert.Equal(t, RecipesPath, cfg.RecipesPath())
	assert.Equal(t, rp, cfg.RecipesPathAbs())
	assert.Equal(t, LockFile, cfg.LockFile())
	assert.Equal(t, MaxAllowedWorkers, cfg.MaxAllowedWorkers())
	assert.Equal(t, MaxAllowedQueueWorkers, cfg.MaxAllowedQueueWorkers())
	assert.Equal(t, ServerPort, cfg.Port())
	assert.Equal(t, WorkerOutputDir, cfg.WorkerOutputDir())
	assert.Equal(t, "", cfg.NotificationsEmail())
	assert.Equal(t, "", cfg.SMTPHost())
	assert.Equal(t, "", cfg.SMTPUsername())
	assert.Equal(t, "", cfg.SMTPPassword())
	assert.Equal(t, []string([]string{"localhost", "www.localhost"}), cfg.ServerDomain())
}

func TestBuildGlobalFallback(t *testing.T) {
	cfg := ConfigGlobalFb

	p, _ := filepath.Abs("../tests/configs/config_global_fb.ini")
	rp, _ := filepath.Abs("../tests/etc/leprechaun/recipes")
	assert.Equal(t, p, cfg.Path())
	assert.Equal(t, "../tests/var/log/leprechaun/error.log", cfg.ErrorLog())
	assert.Equal(t, "../tests/var/log/leprechaun/info.log", cfg.InfoLog())
	assert.Equal(t, "../tests/etc/leprechaun/recipes", cfg.RecipesPath())
	assert.Equal(t, rp, cfg.RecipesPathAbs())
	assert.Equal(t, "../tests/var/run/leprechaun/.lock", cfg.LockFile())
	assert.Equal(t, "../tests/var/log/leprechaun/workers.output", cfg.WorkerOutputDir())
	assert.Equal(t, 5, cfg.MaxAllowedWorkers())
	assert.Equal(t, 5, cfg.MaxAllowedQueueWorkers())
	assert.Equal(t, "some@mail.com", cfg.NotificationsEmail())
	assert.Equal(t, "smtp.host.com", cfg.SMTPHost())
	assert.Equal(t, "smtp_user", cfg.SMTPUsername())
	assert.Equal(t, "smtp_pass", cfg.SMTPPassword())
}

func TestBuildWithSettings(t *testing.T) {
	cfg := ConfigWithSettings

	p, _ := filepath.Abs("../tests/configs/config_regular.ini")
	rp, _ := filepath.Abs("../tests/etc/leprechaun/recipes")
	assert.Equal(t, p, cfg.Path())
	assert.Equal(t, "../tests/var/log/leprechaun/error.log", cfg.ErrorLog())
	assert.Equal(t, "../tests/var/log/leprechaun/info.log", cfg.InfoLog())
	assert.Equal(t, "../tests/etc/leprechaun/recipes", cfg.RecipesPath())
	assert.Equal(t, rp, cfg.RecipesPathAbs())
	assert.Equal(t, "../tests/var/run/leprechaun/.lock", cfg.LockFile())
	assert.Equal(t, "../tests/var/log/leprechaun/workers.output", cfg.WorkerOutputDir())
	assert.Equal(t, 5, cfg.MaxAllowedWorkers())
	assert.Equal(t, 5, cfg.MaxAllowedQueueWorkers())
	assert.Equal(t, "some@mail.com", cfg.NotificationsEmail())
	assert.Equal(t, "smtp.host.com", cfg.SMTPHost())
	assert.Equal(t, "smtp_user", cfg.SMTPUsername())
	assert.Equal(t, "smtp_pass", cfg.SMTPPassword())
	assert.Equal(t, []string{"example.com", "www.example.com"}, cfg.ServerDomain())
}

func TestBuildWithInvalidValues(t *testing.T) {
	cfg := ConfigWithInvalidValues

	assert.Equal(t, ErrorLog, cfg.ErrorLog())
	assert.Equal(t, InfoLog, cfg.InfoLog())
	assert.Equal(t, RecipesPath, cfg.RecipesPath())
	assert.Equal(t, LockFile, cfg.LockFile())
	assert.Equal(t, MaxAllowedWorkers, cfg.MaxAllowedWorkers())
	assert.Equal(t, MaxAllowedQueueWorkers, cfg.MaxAllowedQueueWorkers())
	assert.Equal(t, []string([]string{"localhost", "www.localhost"}), cfg.ServerDomain())
}

func TestGettingNotExistingConfig(t *testing.T) {
	// this should not break, instead we need to get empty object
	configs.Config("not_exists")
}

func TestGettingExistingConfig(t *testing.T) {
	configs.Config("test")
}

func TestNotValidPathToConfig(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			//
		}
	}()

	NewConfigs().New("test", "some_path")
}
