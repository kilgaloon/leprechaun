package config

import (
	"path/filepath"
	"strings"

	"gopkg.in/ini.v1"
)

// Default paths
const (
	ErrorLog               = "/var/log/leprechaun/error.log"
	InfoLog                = "/var/log/leprechaun/info.log"
	RecipesPath            = "/etc/leprechaun/recipes"
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
	ServerDomain           = "localhost"
	ErrorReporting         = true
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

// Config return config by name of the agent
func (c *Configs) Config(name string) *AgentConfig {
	if cfg, ok := c.cfgs[name]; ok {
		return cfg
	}

	return &AgentConfig{}
}

// AgentConfig holds config for agents
type AgentConfig struct {
	path                   string
	errorLog               string
	infoLog                string
	recipesPath            string
	lockFile               string
	workerOutputDir        string
	port                   int
	maxAllowedWorkers      int
	maxAllowedQueueWorkers int
	notificationsEmail     string
	sMTPHost               string
	sMTPUsername           string
	sMTPPassword           string
	domain                 string
	errorReporting         bool
}

// Path returns path of config file
func (ac AgentConfig) Path() string {
	p, err := filepath.Abs(ac.path)
	if err != nil {
		return ac.path
	}

	return p
}

// ErrorLog returns path of config file
func (ac AgentConfig) ErrorLog() string {
	return ac.errorLog
}

// InfoLog returns path of config file
func (ac AgentConfig) InfoLog() string {
	return ac.infoLog
}

// RecipesPathAbs returns absolute path of recipes
func (ac AgentConfig) RecipesPathAbs() string {
	p, err := filepath.Abs(ac.recipesPath)
	if err != nil {
		return ac.recipesPath
	}

	return p
}

// RecipesPath returns path of config file
func (ac AgentConfig) RecipesPath() string {
	return ac.recipesPath
}

// LockFile returns path of config file
func (ac AgentConfig) LockFile() string {
	return ac.lockFile
}

// Port returns path of config file
func (ac AgentConfig) Port() int {
	return ac.port
}

// MaxAllowedWorkers defines how much workers are allowed to work in parallel
func (ac AgentConfig) MaxAllowedWorkers() int {
	return ac.maxAllowedWorkers
}

// MaxAllowedQueueWorkers defines how much workers are allowed to sit in queue
func (ac AgentConfig) MaxAllowedQueueWorkers() int {
	return ac.maxAllowedQueueWorkers
}

// WorkerOutputDir returns path of workers output dir
func (ac AgentConfig) WorkerOutputDir() string {
	return ac.workerOutputDir
}

//NotificationsEmail returns path of workers output dir
func (ac AgentConfig) NotificationsEmail() string {
	return ac.notificationsEmail
}

// SMTPHost returns path of workers output dir
func (ac AgentConfig) SMTPHost() string {
	return ac.sMTPHost
}

// SMTPUsername returns path of workers output dir
func (ac AgentConfig) SMTPUsername() string {
	return ac.sMTPUsername
}

// SMTPPassword returns path of workers output dir
func (ac AgentConfig) SMTPPassword() string {
	return ac.sMTPPassword
}

// ServerDomain returns domain of server
func (ac AgentConfig) ServerDomain() []string {
	var domain string
	var wwwdomain string
	var d []string
	if strings.Contains(ac.domain, "://") {
		d = strings.Split(ac.domain, "://")
	} else {
		d = []string{ac.domain}
	}

	if len(d) > 1 {
		domain = d[1]
	} else {
		domain = d[0]
	}

	if strings.Contains(domain, "www.") {
		wwwdomain = domain
		domain = strings.Trim(domain, "www.")
	} else {
		wwwdomain = "www." + domain
	}

	return []string{domain, wwwdomain}
}

// ErrorReporting returns flag to decide is remote reporting enabled or not
func (ac AgentConfig) ErrorReporting() bool {
	return ac.errorReporting
}

// New Create new config
func (c *Configs) New(name string, path string) *AgentConfig {
	cfg, err := ini.Load(path)
	if err != nil {
		panic(err)
	}

	ac := &AgentConfig{}
	ac.path = path
	gErrorLog := cfg.Section("").Key("error_log").MustString(ErrorLog)
	ac.errorLog = cfg.Section("").Key(name + ".error_log").MustString(gErrorLog)
	if !IsFileValid(ac.errorLog, ".log") {
		ac.errorLog = ErrorLog
	}

	gInfoLog := cfg.Section("").Key("info_log").MustString(InfoLog)
	ac.infoLog = cfg.Section("").Key(name + ".info_log").MustString(gInfoLog)
	if !IsFileValid(ac.infoLog, ".log") {
		ac.infoLog = InfoLog
	}

	gRecipesPath := cfg.Section("").Key("recipes_path").MustString(RecipesPath)
	ac.recipesPath = cfg.Section("").Key(name + ".recipes_path").MustString(gRecipesPath)
	if !IsDirValid(ac.recipesPath) {
		ac.recipesPath = RecipesPath
	}

	gWorkerOutputDir := cfg.Section("").Key("worker_output_dir").MustString(WorkerOutputDir)
	ac.workerOutputDir = cfg.Section("").Key(name + ".worker_output_dir").MustString(gWorkerOutputDir)
	if !IsDirValid(ac.workerOutputDir) {
		ac.workerOutputDir = WorkerOutputDir
	}

	ac.lockFile = cfg.Section("").Key(name + ".lock_file").MustString(LockFile)
	if !IsFileValid(ac.lockFile, ".lock") {
		ac.lockFile = LockFile
	}

	gMaxAllowedWorkers := cfg.Section("").Key("max_allowed_workers").MustInt(MaxAllowedWorkers)
	ac.maxAllowedWorkers = cfg.Section("").Key(name + ".max_allowed_workers").MustInt(gMaxAllowedWorkers)

	gMaxAllowedQueueWorkers := cfg.Section("").Key("max_allowed_queue_workers").MustInt(MaxAllowedQueueWorkers)
	ac.maxAllowedQueueWorkers = cfg.Section("").Key(name + ".max_allowed_queue_workers").MustInt(gMaxAllowedQueueWorkers)

	ac.port = cfg.Section("").Key(name + ".port").MustInt(ServerPort)

	gNotificationsEmail := cfg.Section("").Key("notifications_email").MustString(NotificationsEmail)
	ac.notificationsEmail = cfg.Section("").Key(name + ".notifications_email").MustString(gNotificationsEmail)

	gSMTPHost := cfg.Section("").Key("smtp_host").MustString(SMTPHost)
	ac.sMTPHost = cfg.Section("").Key(name + ".smtp_host").MustString(gSMTPHost)

	gSMTPUsername := cfg.Section("").Key("smtp_username").MustString(SMTPUsername)
	ac.sMTPUsername = cfg.Section("").Key(name + ".smtp_username").MustString(gSMTPUsername)

	gSMTPPassword := cfg.Section("").Key("smtp_password").MustString(SMTPPassword)
	ac.sMTPPassword = cfg.Section("").Key(name + ".smtp_password").MustString(gSMTPPassword)

	ac.domain = cfg.Section("").Key(name + ".domain").MustString(ServerDomain)

	gErrorReporting := cfg.Section("").Key("error_reporting").MustBool(ErrorReporting)
	ac.errorReporting = cfg.Section("").Key(name + ".error_reporting").MustBool(gErrorReporting)

	c.cfgs[name] = ac
	return ac
}
