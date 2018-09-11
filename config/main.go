package config

import (
	"gopkg.in/ini.v1"
)

// Default paths
const (
	// Client
	clientErrorLog    = "/var/log/leprechaun/error.log"
	clientInfoLog     = "/var/log/leprechaun/info.log"
	clientRecipesPath = "/etc/leprechaun/recipes"
	clientPIDFile     = "/var/run/leprechaun/client.pid"
	clientLockFile    = "/var/run/leprechaun/client.lock"

	// Server
	serverErrorLog    = "/var/log/leprechaun/server/error.log"
	serverInfoLog     = "/var/log/leprechaun/server/info.log"
	serverRecipesPath = "/etc/leprechaun/recipes"
	serverPIDFile     = "/var/run/leprechaun/server.pid"
	serverLockFile    = "/var/run/leprechaun/server.lock"
	serverPort        = 11400
)

// Config values
type Config struct {
	ClientConfig
	ServerConfig
}

// ClientConfig holds config for client
type ClientConfig struct {
	ErrorLog    string
	InfoLog     string
	RecipesPath string
	PIDFile     string
	LockFile    string
}

// ServerConfig holds config for server
type ServerConfig struct {
	ErrorLog    string
	InfoLog     string
	RecipesPath string
	PIDFile     string
	LockFile    string
	Port        int
}

// BuildConfig Create client config
func BuildConfig(path string) *Config {
	cfg, err := ini.Load(path)
	if err != nil {
		panic(err)
	}

	c := &Config{}
	c.ClientConfig.ErrorLog = cfg.Section("").Key("client.error_log").MustString(clientErrorLog)
	c.ClientConfig.InfoLog = cfg.Section("").Key("client.info_log").MustString(clientInfoLog)
	c.ClientConfig.RecipesPath = cfg.Section("").Key("client.recipes_path").MustString(clientRecipesPath)
	c.ClientConfig.PIDFile = cfg.Section("").Key("client.pid_file").MustString(clientPIDFile)
	c.ClientConfig.LockFile = cfg.Section("").Key("client.lock_file").MustString(clientLockFile)

	c.ServerConfig.Port = cfg.Section("").Key("server.port").MustInt(serverPort)
	c.ServerConfig.ErrorLog = cfg.Section("").Key("server.error_log").MustString(serverErrorLog)
	c.ServerConfig.InfoLog = cfg.Section("").Key("server.info_log").MustString(serverInfoLog)
	c.ServerConfig.RecipesPath = cfg.Section("").Key("server.recipes_path").MustString(serverRecipesPath)
	c.ServerConfig.PIDFile = cfg.Section("").Key("server.pid_file").MustString(serverPIDFile)
	c.ServerConfig.LockFile = cfg.Section("").Key("server.lock_file").MustString(serverLockFile)

	return c
}

// GetClientConfig returns configuration for client
func (config *Config) GetClientConfig() *ClientConfig {
	return &config.ClientConfig
}

// GetServerConfig returns configuration for server
func (config *Config) GetServerConfig() *ServerConfig {
	return &config.ServerConfig
}
