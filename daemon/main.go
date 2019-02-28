package daemon

import (
	"flag"
	"log"
	"os"
	"strconv"

	"github.com/kilgaloon/leprechaun/api"
	"github.com/kilgaloon/leprechaun/config"
)

// Daemon is long living process that serves as middleware
// and access to multiple agents
type Daemon struct {
	PID        int
	PidPath    string
	PidFile    *os.File
	services   map[string]Service
	Configs    *config.Configs
	ConfigPath string
	Pipeline   chan string
	Cmd        api.Cmd
	Debug      bool
	API        *api.API
}

// Srv is long living process that manages other clients
var Srv *Daemon

// GetPID gets current PID of client
func (d *Daemon) GetPID() int {
	return d.PID
}

// AddService push agent as a service to list of services
func (d *Daemon) AddService(s Service) {
	name := s.GetName()
	cfg := d.Configs.New(name, d.ConfigPath)
	a := s.New(name, cfg, d.Debug)
	a.SetPipeline(d.Pipeline)

	d.API.Register(a)

	d.services[name] = a
}

// Run starts daemon and long living process
func (d *Daemon) Run() {
	if api.IsAPIRunning() {
		api.Resolver(d.Cmd)
	} else {
		go func() {
			for _, s := range d.services {
				log.Printf("Starting service %s", s.GetName())
				go s.Start()
			}

			d.API.RegisterHandle("daemon", d.daemonInfo)
			d.API.Start()
		}()

		for {
			select {
			case info := <-d.Pipeline:
				log.Print(info)
			}
		}

	}

}

// Kill daemon and remove .pid file
func (d *Daemon) Kill() {
	err := os.Remove(d.PidPath)
	if err != nil {
		panic(err)
	}
}

func init() {
	var configPath, pidPath *string
	var debug *bool
	var pid int

	if api.IsAPIRunning() {
		resp := Srv.GetInfo()

		configPath = &resp.ConfigPath
		pidPath = &resp.PidPath
		debug = &resp.Debug
		pid = resp.PID
	} else {
		configPath = flag.String("ini", "../../dist/configs/config.ini", "Path to .ini configuration")
		pidPath = flag.String("pid", "../../tests/var/run/leprechaun/.pid", "PID file of process")
		debug = flag.Bool("debug", false, "Debug mode")
	}

	cmd := flag.String("cmd", "run", "Start all services")
	flag.Parse()

	d := new(Daemon)
	f, err := os.OpenFile(*pidPath, os.O_RDWR|os.O_CREATE, 0644)
	d.PidFile = f
	d.PidPath = *pidPath
	if err != nil {
		panic("Failed to start client, can't save PID, reason: " + err.Error())
	}

	if pid == 0 {
		d.PID = os.Getpid()
		pid := strconv.Itoa(d.PID)
		_, err = d.PidFile.WriteString(pid)
		if err != nil {
			panic("Failed to start client, can't save PID")
		}
	}

	d.services = make(map[string]Service)
	d.ConfigPath = *configPath
	d.Configs = config.NewConfigs()
	d.Pipeline = make(chan string)
	d.Debug = *debug
	d.Cmd = api.Cmd(*cmd)
	d.API = api.New()

	Srv = d
}

//}
