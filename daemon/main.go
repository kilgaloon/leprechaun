package daemon

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

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
	Pipeline     chan string
	Cmd          api.Cmd
	Debug        bool
	API          *api.API
	shutdownChan chan bool
}

// Srv is long living process that manages other clients
var Srv *Daemon

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
	a.SetPipeline(d.Pipeline)

	d.API.Register(a)

	d.services[name] = a
}

// Run starts daemon and long living process
func (d *Daemon) Run() {
	if api.IsAPIRunning() {
		// more commands can/will be used here
		if d.Cmd.Agent() == "daemon" {
			switch d.Cmd.Command() {
			case "info":
				d.renderInfo()
				break
			case "kill":
				d.killDaemon()
				break
			}

			return
		}

		api.Resolver(d.Cmd)
	} else {
		go func() {
			for _, s := range d.services {
				log.Printf("Starting service %s", s.GetName())
				go s.Start()
			}

			d.API.RegisterHandle("daemon/info", d.daemonInfo)
			d.API.RegisterHandle("daemon/kill", d.daemonKill)
			d.API.Start()
		}()

		for {
			select {
			case info := <-d.Pipeline:
				log.Print(info)
			case <-d.shutdownChan:
				if os.Getenv("RUN_MODE") != "test" {
					os.Exit(1)
				}
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

func init() {
	var configPath, pidPath *string
	var debug, helpFlag *bool
	var pid int

	helpFlag = flag.Bool("help", false, "Display helpfull info")

	if api.IsAPIRunning() {
		resp := Srv.GetInfo()

		configPath = &resp.ConfigPath
		pidPath = &resp.PidPath
		debug = &resp.Debug
		pid = resp.PID
	} else {
		if os.Getenv("RUN_MODE") == "test" {
			pp := "../tests/var/run/leprechaun/.pid"
			cp := "../tests/configs/config_regular.ini"
			dbg := true

			pidPath = &pp
			configPath = &cp
			debug = &dbg
		} else {
			configPath = flag.String("ini", "/etc/leprechaun/config.ini", "Path to .ini configuration")
			pidPath = flag.String("pid", "/var/run/leprechaun/.pid", "PID file of process")
			debug = flag.Bool("debug", false, "Debug mode")
		}

	}

	cmd := flag.String("cmd", "run", "Send commands to agents and they will respond")
	flag.Parse()

	if *helpFlag {
		help := "\nAvailable commands for leprechaun --cmd='{agent} {command} {args}' \n" +
			"====== \n" +
			"daemon info - Display basic informations about daemon. \n" +
			"daemon kill - Kills process. \n" +
			"====== \n" +
			"{agent} info - Display basic info about agent.\n" +
			"{agent} start - Start agent if its stopped/paused.\n" +
			"{agent} stop - Stop agent, note that this will remove everything from memory and starting will rebuild agent from scratch.\n" +
			"{agent} pause - Pause agent will not remove everything from memory and if started again it will just continue.\n" +
			"{agent} workers:list - Show list of currently active workers for agent and some basic info.\n" +
			"{agent} workers:kill {name} - Kill worker that match name provided.\n"

		fmt.Println(help)

		os.Exit(1)
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
	d.Pipeline = make(chan string)
	d.Debug = *debug
	d.Cmd = api.Cmd(*cmd)
	d.API = api.New()
	d.shutdownChan = make(chan bool, 1)

	Srv = d
}
