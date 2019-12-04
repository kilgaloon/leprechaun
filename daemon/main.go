package daemon

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/getsentry/raven-go"
	"github.com/kilgaloon/leprechaun/api"
	"github.com/kilgaloon/leprechaun/config"
)

// Daemon is long living process that serves as middleware
// and access to multiple agents
type Daemon struct {
	PID          int
	PidPath      string
	PidFile      *os.File
	services     map[string]Service
	Configs      *config.Configs
	ConfigPath   string
	Cmd          api.Cmd
	Debug        bool
	API          *api.API
	shutdownChan chan bool
}

// GetPID gets current PID of client
func (d *Daemon) GetPID() int {
	return d.PID
}

// GetConfigPath returns path of config file
func (d *Daemon) GetConfigPath() string {
	p, err := filepath.Abs(d.ConfigPath)
	if err != nil {
		return d.ConfigPath
	}

	return p
}

// GetPidPath returns path of config file
func (d *Daemon) GetPidPath() string {
	p, err := filepath.Abs(d.PidPath)
	if err != nil {
		return d.PidPath
	}

	return p
}

// AddService push agent as a service to list of services
func (d *Daemon) AddService(s Service) {
	name := s.GetName()
	cfg := d.Configs.New(name, d.GetConfigPath())
	a := s.New(name, cfg, d.Debug)

	d.API.Register(a)

	d.services[name] = a
}

// Run starts daemon and long living process
func (d *Daemon) Run(cb func()) {
	if api.IsAPIRunning() {
		// more commands can/will be used here
		if d.Cmd.Agent() == "daemon" {
			switch d.Cmd.Command() {
			case "info":
				d.renderInfo()
				return
			}
		}

		api.Resolver(d.Cmd)
	} else {
		var servicesToStart []string
		if d.Cmd.Agent() == "run" {
			servicesToStart = strings.Split(d.Cmd.Command(), ",")
		}

		for _, s := range d.services {
			for _, srv := range servicesToStart {
				if srv == s.GetName() {
					log.Printf("Starting service %s", s.GetName())
					go s.Start()

					break
				}
			}
		}

		d.API.RegisterHandle("daemon/info", d.daemonInfo)
		d.API.RegisterHandle("daemon/kill", d.daemonKill)
		d.API.RegisterHandle("daemon/services", d.servicesList)
		d.API.Start()

		if cb != nil {
			cb()
		}

		for {
			select {
			case <-d.shutdownChan:
				if os.Getenv("RUN_MODE") != "test" {
					os.Exit(1)
				}

				d.API.Stop()

				break
			}
		}

	}
}

// Kill daemon and remove .pid file
func (d *Daemon) Kill() {
	err := os.Remove(d.GetPidPath())
	if err != nil {
		panic(err)
	}

	d.shutdownChan <- true
}

//Init initialize daemon
func Init() *Daemon {
	var configPath, pidPath, cmd *string
	var debug, hf *bool
	var pid int

	if api.IsAPIRunning() {
		d := &Daemon{}
		resp := d.GetInfo()

		configPath = &resp.ConfigPath
		pidPath = &resp.PidPath
		debug = &resp.Debug
		pid = resp.PID
	} else {
		if !flag.Parsed() {
			hf = flag.Bool("commands", false, "Display helpful info")
			cmd = flag.String("cmd", "run scheduler,server,cron", "Send commands to agents and they will respond (default command is to run all services)")
			configPath = flag.String("ini", "/etc/leprechaun/config.ini", "Path to .ini configuration")
			pidPath = flag.String("pid", "/var/run/leprechaun/.pid", "PID file of process")
			debug = flag.Bool("debug", false, "Debug mode")

			flag.Parse()
		}
	}

	if hf != nil {
		if *hf {
			helpCommands()

			return nil
		}
	}

	if os.Getenv("RUN_MODE") == "test" {
		pp := "../tests/var/run/leprechaun/.pid"
		cp := "../tests/configs/config_regular.ini"
		dbg := true
		cmd = new(string)

		pidPath = &pp
		configPath = &cp
		debug = &dbg
	}

	d := new(Daemon)
	f, err := os.OpenFile(*pidPath, os.O_RDWR|os.O_CREATE, 0644)
	d.PidFile = f
	d.PidPath = *pidPath
	if err != nil {
		log.Fatal("Failed to start client, can't save PID. Directory for pid file doesn't exist or pid file not valid")
	}

	if pid == 0 {
		d.PID = os.Getpid()
		pid := strconv.Itoa(d.PID)
		_, err = d.PidFile.WriteString(pid)
		if err != nil {
			log.Fatal("Failed to start client, can't save PID")
		}
	}

	d.services = make(map[string]Service)
	d.ConfigPath = *configPath
	d.Configs = config.NewConfigs()
	d.Debug = *debug
	d.Cmd = api.Cmd(*cmd)
	d.API = api.New()
	d.shutdownChan = make(chan bool, 1)

	cfg := d.Configs.New("daemon", d.ConfigPath)
	if cfg.GetErrorReporting() && os.Getenv("RUN_MODE") != "test" {
		raven.SetDSN("https://63f6916e9a4f4ae08853f5b1fe5eabda:a2abd1e7a4f944dca632f875279197f7@sentry.io/1422644")
	}

	return d
}
