package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	ConfigWithoutDefaultSettings = "../tests/configs/config_without_default_values.ini"
	ConfigWithSettings           = "../tests/configs/config_regular.ini"
	ConfigWithInvalidValues      = "../tests/configs/config_wrong_value.ini"
	ConfigWithWrongExt           = "../tests/configs/config_wrong_ext.ini"
)

func TestBuildWithoutSettings(t *testing.T) {
	configs := NewConfigs()
	configs.New("test", ConfigWithoutDefaultSettings)
	cfg := configs.GetConfig("test")

	assert.Equal(t, ConfigWithoutDefaultSettings, cfg.GetPath())
	assert.Equal(t, ErrorLog, cfg.GetErrorLog())
	assert.Equal(t, InfoLog, cfg.GetInfoLog())
	assert.Equal(t, RecipesPath, cfg.GetRecipesPath())
	assert.Equal(t, PIDFile, cfg.GetPIDFile())
	assert.Equal(t, LockFile, cfg.GetLockFile())
	assert.Equal(t, MaxAllowedWorkers, cfg.GetMaxAllowedWorkers())
	assert.Equal(t, RetryRecipeAfter, cfg.GetRetryRecipeAfter())
	assert.Equal(t, CommandSocket, cfg.GetCommandSocket())
	assert.Equal(t, ServerPort, cfg.GetPort())
}

func TestBuildWithSettings(t *testing.T) {
	cfg := NewConfigs().New("test", ConfigWithSettings)

	assert.Equal(t, "../tests/var/log/leprechaun/client-error.log", cfg.GetErrorLog())
	assert.Equal(t, "../tests/var/log/leprechaun/client-info.log", cfg.GetInfoLog())
	assert.Equal(t, "../tests/etc/leprechaun/recipes", cfg.GetRecipesPath())
	assert.Equal(t, "../tests/var/run/leprechaun/client.pid", cfg.GetPIDFile())
	assert.Equal(t, "../tests/var/run/leprechaun/client.lock", cfg.GetLockFile())
	assert.Equal(t, 5, cfg.GetMaxAllowedWorkers())
	assert.Equal(t, 10, cfg.GetRetryRecipeAfter())
}

func TestBuildWithSettingsWithWrongExt(t *testing.T) {
	NewConfigs().New("test", ConfigWithWrongExt)
}

func TestBuildWithInvalidValues(t *testing.T) {
	cfg := NewConfigs().New("test", ConfigWithInvalidValues)

	assert.Equal(t, ErrorLog, cfg.GetErrorLog())
	assert.Equal(t, InfoLog, cfg.GetInfoLog())
	assert.Equal(t, RecipesPath, cfg.GetRecipesPath())
	assert.Equal(t, PIDFile, cfg.GetPIDFile())
	assert.Equal(t, LockFile, cfg.GetLockFile())
	assert.Equal(t, MaxAllowedWorkers, cfg.GetMaxAllowedWorkers())
	assert.Equal(t, RetryRecipeAfter, cfg.GetRetryRecipeAfter())
}

func TestNotValidPathToConfig(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			//
		}
	}()

	NewConfigs().New("test", "some_path")
}
