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
	assert.Equal(t, clientCfg.ErrorLog, clientErrorLog)
	assert.Equal(t, clientCfg.InfoLog, clientInfoLog)
	assert.Equal(t, clientCfg.RecipesPath, clientRecipesPath)
	assert.Equal(t, clientCfg.PIDFile, clientPIDFile)
	assert.Equal(t, clientCfg.LockFile, clientLockFile)

	serverCfg := cfg.GetServerConfig()
	assert.Equal(t, serverCfg.ErrorLog, serverErrorLog)
	assert.Equal(t, serverCfg.InfoLog, serverInfoLog)
	assert.Equal(t, serverCfg.RecipesPath, serverRecipesPath)
	assert.Equal(t, serverCfg.Port, serverPort)
	assert.Equal(t, serverCfg.PIDFile, serverPIDFile)
	assert.Equal(t, serverCfg.LockFile, serverLockFile)
}

func TestBuildWithSettings(t *testing.T) {
	cfg := BuildConfig(ConfigWithSettings)

	clientCfg := cfg.GetClientConfig()
	assert.Equal(t, clientCfg.ErrorLog, "../tests/var/log/leprechaun/client-error.log")
	assert.Equal(t, clientCfg.InfoLog, "../tests/var/log/leprechaun/client-info.log")
	assert.Equal(t, clientCfg.RecipesPath, "../tests/etc/leprechaun/recipes")
	assert.Equal(t, clientCfg.PIDFile, "../tests/var/run/leprechaun/client.pid")
	assert.Equal(t, clientCfg.LockFile, "../tests/var/run/leprechaun/client.lock")

	serverCfg := cfg.GetServerConfig()
	assert.Equal(t, serverCfg.ErrorLog, "../tests/var/log/leprechaun/server-error.log")
	assert.Equal(t, serverCfg.InfoLog, "../tests/var/log/leprechaun/server-info.log")
	assert.Equal(t, serverCfg.RecipesPath, "../tests/etc/leprechaun/recipes")
	assert.Equal(t, serverCfg.Port, 11400)
	assert.Equal(t, serverCfg.PIDFile, "../tests/var/run/leprechaun/server.pid")
	assert.Equal(t, serverCfg.LockFile, "../tests/var/run/leprechaun/server.lock")
}

func TestBuildWithInvalidValues(t *testing.T) {
	cfg := BuildConfig(ConfigWithoutDefaultSettings)

	clientCfg := cfg.GetClientConfig()
	assert.Equal(t, clientCfg.ErrorLog, clientErrorLog)
	assert.Equal(t, clientCfg.InfoLog, clientInfoLog)
	assert.Equal(t, clientCfg.RecipesPath, clientRecipesPath)
	assert.Equal(t, clientCfg.PIDFile, clientPIDFile)
	assert.Equal(t, clientCfg.LockFile, clientLockFile)

	serverCfg := cfg.GetServerConfig()
	assert.Equal(t, serverCfg.ErrorLog, serverErrorLog)
	assert.Equal(t, serverCfg.InfoLog, serverInfoLog)
	assert.Equal(t, serverCfg.RecipesPath, serverRecipesPath)
	assert.Equal(t, serverCfg.Port, serverPort)
	assert.Equal(t, serverCfg.PIDFile, serverPIDFile)
	assert.Equal(t, serverCfg.LockFile, serverLockFile)
}
