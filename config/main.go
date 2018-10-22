package config

import (
	"gopkg.in/ini.v1"
)

// Default paths
const (
	ErrorLog               = "/var/log/leprechaun/error.log"
	InfoLog                = "/var/log/leprechaun/info.log"
	RecipesPath            = "/etc/leprechaun/recipes"
	PIDFile                = "/var/run/leprechaun/client.pid"
	LockFile               = "/var/run/leprechaun/client.lock"
	CommandSocket          = "/var/run/leprechaun/client.sock"
	WorkerOutputDir        = "/var/log/leprechaun/workers.output"
	NotificationsEmail     = ""
	MaxAllowedWorkers      = 5
	MaxAllowedQueueWorkers = 5
	RetryRecipeAfter       = 10
	ServerPort             = 11400
	SMTPHost               = ""
	SMTPUsername           = ""
	SMTPPassword           = ""
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
	Path                   string
	ErrorLog               string
	InfoLog                string
	RecipesPath            string
	PIDFile                string
	LockFile               string
	CommandSocket          string
	WorkerOutputDir        string
	Port                   int
	MaxAllowedWorkers      int
	MaxAllowedQueueWorkers int
	NotificationsEmail     string
	SMTPHost               string
	SMTPUsername           string
	SMTPPassword           string
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

// GetMaxAllowedWorkers defines how much workers are allowed to work in parallel
func (ac AgentConfig) GetMaxAllowedWorkers() int {
	return ac.MaxAllowedWorkers
}

// GetMaxAllowedQueueWorkers defines how much workers are allowed to sit in queue
func (ac AgentConfig) GetMaxAllowedQueueWorkers() int {
	return ac.MaxAllowedQueueWorkers
}

// GetWorkerOutputDir returns path of workers output dir
func (ac AgentConfig) GetWorkerOutputDir() string {
	return ac.WorkerOutputDir
}

// GetNotificationsEmail returns path of workers output dir
func (ac AgentConfig) GetNotificationsEmail() string {
	return ac.NotificationsEmail
}

// GetSMTPHost returns path of workers output dir
func (ac AgentConfig) GetSMTPHost() string {
	return ac.SMTPHost
}

// GetSMTPUsername returns path of workers output dir
func (ac AgentConfig) GetSMTPUsername() string {
	return ac.SMTPUsername
}

// GetSMTPPassword returns path of workers output dir
func (ac AgentConfig) GetSMTPPassword() string {
	return ac.SMTPPassword
}

// New Create new config
func (c *Configs) New(name string, path string) *AgentConfig {
	cfg, err := ini.Load(path)
	if err != nil {
		panic(err)
	}

	ac := &AgentConfig{}
	ac.Path = path
	gErrorLog := cfg.Section("").Key("error_log").MustString(ErrorLog)
	ac.ErrorLog = cfg.Section("").Key(name + ".error_log").MustString(gErrorLog)
	if !IsFileValid(ac.ErrorLog, ".log") {
		ac.ErrorLog = ErrorLog
	}

	gInfoLog := cfg.Section("").Key("info_log").MustString(InfoLog)
	ac.InfoLog = cfg.Section("").Key(name + ".info_log").MustString(gInfoLog)
	if !IsFileValid(ac.InfoLog, ".log") {
		ac.InfoLog = InfoLog
	}

	gRecipesPath := cfg.Section("").Key("recipes_path").MustString(RecipesPath)
	ac.RecipesPath = cfg.Section("").Key(name + ".recipes_path").MustString(gRecipesPath)
	if !IsDirValid(ac.RecipesPath) {
		ac.RecipesPath = RecipesPath
	}

	gWorkerOutputDir := cfg.Section("").Key("worker_output_dir").MustString(WorkerOutputDir)
	ac.WorkerOutputDir = cfg.Section("").Key(name + ".worker_output_dir").MustString(gWorkerOutputDir)
	if !IsDirValid(ac.WorkerOutputDir) {
		ac.WorkerOutputDir = WorkerOutputDir
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

	gMaxAllowedWorkers := cfg.Section("").Key("max_allowed_workers").MustInt(MaxAllowedWorkers)
	ac.MaxAllowedWorkers = cfg.Section("").Key(name + ".max_allowed_workers").MustInt(gMaxAllowedWorkers)

	gMaxAllowedQueueWorkers := cfg.Section("").Key("max_allowed_queue_workers").MustInt(MaxAllowedQueueWorkers)
	ac.MaxAllowedQueueWorkers = cfg.Section("").Key(name + ".max_allowed_queue_workers").MustInt(gMaxAllowedQueueWorkers)

	ac.Port = cfg.Section("").Key(name + ".port").MustInt(ServerPort)

	gNotificationsEmail := cfg.Section("").Key("notifications_email").MustString(NotificationsEmail)
	ac.NotificationsEmail = cfg.Section("").Key(name + ".notifications_email").MustString(gNotificationsEmail)

	gSMTPHost := cfg.Section("").Key("smtp_host").MustString(SMTPHost)
	ac.SMTPHost = cfg.Section("").Key(name + ".smtp_host").MustString(gSMTPHost)

	gSMTPUsername := cfg.Section("").Key(name + ".smtp_username").MustString(SMTPUsername)
	ac.SMTPUsername = cfg.Section("").Key(name + ".smtp_username").MustString(gSMTPUsername)

	gSMTPPassword := cfg.Section("").Key("smtp_password").MustString(SMTPPassword)
	ac.SMTPPassword = cfg.Section("").Key(name + ".smtp_password").MustString(gSMTPPassword)

	c.cfgs[name] = ac
	return ac
}
