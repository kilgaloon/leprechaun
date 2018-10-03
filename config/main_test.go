package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	configs                      = NewConfigs()
	ConfigWithoutDefaultSettings = configs.New("test", "../tests/configs/config_without_default_values.ini")
	ConfigWithSettings           = configs.New("test", "../tests/configs/config_regular.ini")
	ConfigWithInvalidValues      = configs.New("test", "../tests/configs/config_wrong_value.ini")
	ConfigWithWrongExt           = configs.New("test", "../tests/configs/config_wrong_ext.ini")
)

func TestBuildWithoutSettings(t *testing.T) {
	cfg := ConfigWithoutDefaultSettings

	assert.Equal(t, "../tests/configs/config_without_default_values.ini", cfg.GetPath())
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
	cfg := ConfigWithSettings

	assert.Equal(t, "../tests/configs/config_regular.ini", cfg.GetPath())
	assert.Equal(t, "../tests/var/log/leprechaun/error.log", cfg.GetErrorLog())
	assert.Equal(t, "../tests/var/log/leprechaun/info.log", cfg.GetInfoLog())
	assert.Equal(t, "../tests/etc/leprechaun/recipes", cfg.GetRecipesPath())
	assert.Equal(t, "../tests/var/run/leprechaun/.pid", cfg.GetPIDFile())
	assert.Equal(t, "../tests/var/run/leprechaun/.lock", cfg.GetLockFile())
	assert.Equal(t, "../tests/var/run/leprechaun/.sock", cfg.GetCommandSocket())
	assert.Equal(t, 5, cfg.GetMaxAllowedWorkers())
	assert.Equal(t, 10, cfg.GetRetryRecipeAfter())
}

func TestBuildWithInvalidValues(t *testing.T) {
	cfg := ConfigWithInvalidValues

	assert.Equal(t, ErrorLog, cfg.GetErrorLog())
	assert.Equal(t, InfoLog, cfg.GetInfoLog())
	assert.Equal(t, RecipesPath, cfg.GetRecipesPath())
	assert.Equal(t, PIDFile, cfg.GetPIDFile())
	assert.Equal(t, LockFile, cfg.GetLockFile())
	assert.Equal(t, MaxAllowedWorkers, cfg.GetMaxAllowedWorkers())
	assert.Equal(t, RetryRecipeAfter, cfg.GetRetryRecipeAfter())
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
