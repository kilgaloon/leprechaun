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
	CertPemPath            = "/etc/leprechaun/certs/server.pem"
	CertKeyPath            = "/etc/leprechaun/certs/server.key"
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
	Cfg                    *ini.File
	Path                   string
	ErrorLog               string
	InfoLog                string
	RecipesPath            string
	LockFile               string
	WorkerOutputDir        string
	Port                   int
	MaxAllowedWorkers      int
	MaxAllowedQueueWorkers int
	NotificationsEmail     string
	SMTPHost               string
	SMTPUsername           string
	SMTPPassword           string
	Domain                 string
	ErrorReporting         bool
	CertPemPath            string
	CertKeyPath            string
	RemoteHosts            map[string]string
}

// GetPath returns path of config file
func (ac AgentConfig) GetPath() string {
	p, err := filepath.Abs(ac.Path)
	if err != nil {
		return ac.Path
	}

	return p
}

// GetErrorLog returns path of config file
func (ac AgentConfig) GetErrorLog() string {
	return ac.ErrorLog
}

// GetInfoLog returns path of config file
func (ac AgentConfig) GetInfoLog() string {
	return ac.InfoLog
}

// GetRecipesPathAbs returns absolute path of recipes
func (ac AgentConfig) GetRecipesPathAbs() string {
	p, err := filepath.Abs(ac.RecipesPath)
	if err != nil {
		return ac.RecipesPath
	}

	return p
}

// GetRecipesPath returns path of config file
func (ac AgentConfig) GetRecipesPath() string {
	return ac.RecipesPath
}

// GetLockFile returns path of config file
func (ac AgentConfig) GetLockFile() string {
	return ac.LockFile
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

// GetServerDomain returns domain of server
func (ac AgentConfig) GetServerDomain() []string {
	var domain string
	var wwwdomain string
	var d []string
	if strings.Contains(ac.Domain, "://") {
		d = strings.Split(ac.Domain, "://")
	} else {
		d = []string{ac.Domain}
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

// GetErrorReporting returns flag to decide is remote reporting enabled or not
func (ac AgentConfig) GetErrorReporting() bool {
	return ac.ErrorReporting
}

// GetCertPemPath returns path to .pem file
func (ac AgentConfig) GetCertPemPath() string {
	return ac.CertPemPath
}

// GetCertKeyPath returns path to .key file
func (ac AgentConfig) GetCertKeyPath() string {
	return ac.CertKeyPath
}

// GetRemoteServices map remote services with ports they are running on
func (ac AgentConfig) GetRemoteServices() map[string]string {
	mapped := make(map[string]string)
	keys := ac.Cfg.Section("remote_services").Keys()

	for _, key := range keys {
		mapped[key.Name()] = key.Value()
	}

	return mapped
}

// New Create new config
func (c *Configs) New(name string, path string) *AgentConfig {
	cfg, err := ini.Load(path)
	if err != nil {
		panic(err)
	}

	ac := &AgentConfig{
		Cfg: cfg,
	}

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

	ac.LockFile = cfg.Section("").Key(name + ".lock_file").MustString(LockFile)
	if !IsFileValid(ac.LockFile, ".lock") {
		ac.LockFile = LockFile
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

	gSMTPUsername := cfg.Section("").Key("smtp_username").MustString(SMTPUsername)
	ac.SMTPUsername = cfg.Section("").Key(name + ".smtp_username").MustString(gSMTPUsername)

	gSMTPPassword := cfg.Section("").Key("smtp_password").MustString(SMTPPassword)
	ac.SMTPPassword = cfg.Section("").Key(name + ".smtp_password").MustString(gSMTPPassword)

	ac.Domain = cfg.Section("").Key(name + ".domain").MustString(ServerDomain)

	gErrorReporting := cfg.Section("").Key("error_reporting").MustBool(ErrorReporting)
	ac.ErrorReporting = cfg.Section("").Key(name + ".error_reporting").MustBool(gErrorReporting)

	gCertPemPath := cfg.Section("").Key("pem_file").MustString(CertPemPath)
	ac.CertPemPath = cfg.Section("").Key(name + ".pem_file").MustString(gCertPemPath)

	gCertKeyPath := cfg.Section("").Key("key_file").MustString(CertKeyPath)
	ac.CertKeyPath = cfg.Section("").Key(name + ".key_file").MustString(gCertKeyPath)

	c.cfgs[name] = ac
	return ac
}
