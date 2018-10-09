package config

import (
	"gopkg.in/ini.v1"
)

// Default paths
const (
	ErrorLog          = "/var/log/leprechaun/error.log"
	InfoLog           = "/var/log/leprechaun/info.log"
	RecipesPath       = "/etc/leprechaun/recipes"
	PIDFile           = "/var/run/leprechaun/client.pid"
	LockFile          = "/var/run/leprechaun/client.lock"
	CommandSocket     = "/var/run/leprechaun/client.sock"
	WorkerOutputDir   = "/var/log/leprechaun/workers.output"
	MaxAllowedWorkers = 5
	RetryRecipeAfter  = 10
	ServerPort        = 11400
)

// Configs for different agents
type Configs struct {
	cfgs map[string]*AgentConfig
}

// NewConfigs return Configs of agents
func NewConfigs() *Configs {
	return &Configs{
		cfgs: make(map[string]*AgentConfig),
	}
}

// GetConfig return config by name of the agent
func (c *Configs) GetConfig(name string) *AgentConfig {
	if cfg, ok := c.cfgs[name]; ok {
		return cfg
	}

	return &AgentConfig{}
}

// AgentConfig holds config for agents
type AgentConfig struct {
	Path              string
	ErrorLog          string
	InfoLog           string
	RecipesPath       string
	PIDFile           string
	LockFile          string
	CommandSocket     string
	WorkerOutputDir   string
	Port              int
	MaxAllowedWorkers int
	RetryRecipeAfter  int
}

// GetPath returns path of config file
func (ac AgentConfig) GetPath() string {
	return ac.Path
}

// GetErrorLog returns path of config file
func (ac AgentConfig) GetErrorLog() string {
	return ac.ErrorLog
}

// GetInfoLog returns path of config file
func (ac AgentConfig) GetInfoLog() string {
	return ac.InfoLog
}

// GetRecipesPath returns path of config file
func (ac AgentConfig) GetRecipesPath() string {
	return ac.RecipesPath
}

// GetPIDFile returns path of config file
func (ac AgentConfig) GetPIDFile() string {
	return ac.PIDFile
}

// GetLockFile returns path of config file
func (ac AgentConfig) GetLockFile() string {
	return ac.LockFile
}

// GetCommandSocket returns path of config file
func (ac AgentConfig) GetCommandSocket() string {
	return ac.CommandSocket
}

// GetPort returns path of config file
func (ac AgentConfig) GetPort() int {
	return ac.Port
}

// GetMaxAllowedWorkers returns path of config file
func (ac AgentConfig) GetMaxAllowedWorkers() int {
	return ac.MaxAllowedWorkers
}

// GetRetryRecipeAfter returns path of config file
func (ac AgentConfig) GetRetryRecipeAfter() int {
	return ac.RetryRecipeAfter
}

// GetWorkerOutputDir returns path of workers output dir
func (ac AgentConfig) GetWorkerOutputDir() string {
	return ac.WorkerOutputDir
}

// New Create new config
func (c *Configs) New(name string, path string) *AgentConfig {
	cfg, err := ini.Load(path)
	if err != nil {
		panic(err)
	}

	ac := &AgentConfig{}
	ac.Path = path
	ac.ErrorLog = cfg.Section("").Key(name + ".error_log").MustString(ErrorLog)
	if !IsFileValid(ac.ErrorLog, ".log") {
		ac.ErrorLog = ErrorLog
	}

	ac.InfoLog = cfg.Section("").Key(name + ".info_log").MustString(InfoLog)
	if !IsFileValid(ac.InfoLog, ".log") {
		ac.InfoLog = InfoLog
	}

	ac.RecipesPath = cfg.Section("").Key(name + ".recipes_path").MustString(RecipesPath)
	if !IsDirValid(ac.RecipesPath) {
		ac.RecipesPath = RecipesPath
	}

	ac.WorkerOutputDir = cfg.Section("").Key(name + ".worker_output_dir").MustString(RecipesPath)
	if !IsDirValid(ac.RecipesPath) {
		ac.RecipesPath = RecipesPath
	}

	ac.PIDFile = cfg.Section("").Key(name + ".pid_file").MustString(PIDFile)
	if !IsFileValid(ac.PIDFile, ".pid") {
		ac.PIDFile = PIDFile
	}

	ac.LockFile = cfg.Section("").Key(name + ".lock_file").MustString(LockFile)
	if !IsFileValid(ac.LockFile, ".lock") {
		ac.LockFile = LockFile
	}

	ac.CommandSocket = cfg.Section("").Key(name + ".command_socket").MustString(CommandSocket)
	if !IsFileValid(ac.CommandSocket, ".sock") {
		ac.CommandSocket = CommandSocket
	}

	ac.MaxAllowedWorkers = cfg.Section("").Key(name + ".max_allowed_workers").MustInt(MaxAllowedWorkers)
	ac.RetryRecipeAfter = cfg.Section("").Key(name + ".retry_recipe_after").MustInt(RetryRecipeAfter)
	ac.Port = cfg.Section("").Key(name + ".port").MustInt(ServerPort)

	c.cfgs[name] = ac
	return ac
}
