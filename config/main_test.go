package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	ConfigWithoutDefaultSettings = "../tests/configs/config_without_default_values.ini"
	ConfigWithSettings           = "../tests/configs/config_regular.ini"
	ConfigWithInvalidValues      = "../tests/configs/config_wrong_value.ini"
)

func TestBuildWithoutSettings(t *testing.T) {
	cfg := BuildConfig(ConfigWithoutDefaultSettings)

	clientCfg := cfg.GetClientConfig()
	assert.Equal(t, clientErrorLog, clientCfg.ErrorLog)
	assert.Equal(t, clientInfoLog, clientCfg.InfoLog)
	assert.Equal(t, clientRecipesPath, clientCfg.RecipesPath)
	assert.Equal(t, clientPIDFile, clientCfg.PIDFile)
	assert.Equal(t, clientLockFile, clientCfg.LockFile)

	serverCfg := cfg.GetServerConfig()
	assert.Equal(t, serverErrorLog, serverCfg.ErrorLog)
	assert.Equal(t, serverInfoLog, serverCfg.InfoLog)
	assert.Equal(t, serverRecipesPath, serverCfg.RecipesPath)
	assert.Equal(t, serverPort, serverCfg.Port)
	assert.Equal(t, serverPIDFile, serverCfg.PIDFile)
	assert.Equal(t, serverLockFile, serverCfg.LockFile)
}

func TestBuildWithSettings(t *testing.T) {
	cfg := BuildConfig(ConfigWithSettings)

	clientCfg := cfg.GetClientConfig()
	assert.Equal(t, "../tests/var/log/leprechaun/client-error.log", clientCfg.ErrorLog)
	assert.Equal(t, "../tests/var/log/leprechaun/client-info.log", clientCfg.InfoLog)
	assert.Equal(t, "../tests/etc/leprechaun/recipes", clientCfg.RecipesPath)
	assert.Equal(t, "../tests/var/run/leprechaun/client.pid", clientCfg.PIDFile)
	assert.Equal(t, "../tests/var/run/leprechaun/client.lock", clientCfg.LockFile)

	serverCfg := cfg.GetServerConfig()
	assert.Equal(t, "../tests/var/log/leprechaun/server-error.log", serverCfg.ErrorLog)
	assert.Equal(t, "../tests/var/log/leprechaun/server-info.log", serverCfg.InfoLog)
	assert.Equal(t, "../tests/etc/leprechaun/recipes", serverCfg.RecipesPath)
	assert.Equal(t, 11400, serverCfg.Port)
	assert.Equal(t, "../tests/var/run/leprechaun/server.pid", serverCfg.PIDFile)
	assert.Equal(t, "../tests/var/run/leprechaun/server.lock", serverCfg.LockFile)
}

func TestBuildWithInvalidValues(t *testing.T) {
	cfg := BuildConfig(ConfigWithoutDefaultSettings)

	clientCfg := cfg.GetClientConfig()
	assert.Equal(t, clientErrorLog, clientCfg.ErrorLog)
	assert.Equal(t, clientInfoLog, clientCfg.InfoLog)
	assert.Equal(t, clientRecipesPath, clientCfg.RecipesPath)
	assert.Equal(t, clientPIDFile, clientCfg.PIDFile)
	assert.Equal(t, clientLockFile, clientCfg.LockFile)

	serverCfg := cfg.GetServerConfig()
	assert.Equal(t, serverErrorLog, serverCfg.ErrorLog)
	assert.Equal(t, serverInfoLog, serverCfg.InfoLog)
	assert.Equal(t, serverRecipesPath, serverCfg.RecipesPath)
	assert.Equal(t, serverPort, serverCfg.Port)
	assert.Equal(t, serverPIDFile, serverCfg.PIDFile)
	assert.Equal(t, serverLockFile, serverCfg.LockFile)
}
