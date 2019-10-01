package remote

import (
	"bytes"
	"crypto/rand"
	"crypto/tls"
	"net"
	"net/http"

	"github.com/kilgaloon/leprechaun/agent"
	"github.com/kilgaloon/leprechaun/config"
	"github.com/kilgaloon/leprechaun/daemon"
	"github.com/kilgaloon/leprechaun/workers"
)

// Agent holds client instance
var Agent *Remote

// Remote settings and configurations
type Remote struct {
	Name string
	*agent.Default
}

// New create remote as a service
func (r *Remote) New(name string, cfg *config.AgentConfig, debug bool) daemon.Service {
	a := agent.New(name, cfg, debug)
	c := &Remote{
		name,
		a,
	}

	Agent = c

	return c
}

// GetName returns service name
func (r *Remote) GetName() string {
	return r.Name
}

// Start remote service
func (r *Remote) Start() {
	r.SetStatus(daemon.Started)
	if !r.IsDebug() {
		cert, err := tls.LoadX509KeyPair(r.GetConfig().GetCertPemPath(), r.GetConfig().GetCertKeyPath())
		if err != nil {
			r.Error(err.Error())
		}

		config := tls.Config{Certificates: []tls.Certificate{cert}, ClientAuth: tls.RequireAnyClientCert}
		config.Rand = rand.Reader

		ln, err := tls.Listen("tcp", ":11402", &config)
		if err != nil {
			r.Error(err.Error())
		}

		r.Info("Server(TLS) up and listening on port 11402")

		for {
			conn, err := ln.Accept()
			if err != nil {
				r.Error(err.Error())
				continue
			}
			go r.handleConnection(conn)
		}
	} else {
		ln, err := net.Listen("tcp", ":11402")
		if err != nil {
			r.Error(err.Error())
		}

		for {
			conn, err := ln.Accept()
			if err != nil {
				r.Error(err.Error())
				continue
			}
			go r.handleConnection(conn)
		}
	}

}

func (r *Remote) handleConnection(c net.Conn) {
	r.Info("Client %v connected.", c.RemoteAddr())

	buffer := make([]byte, 256)

	for {
		_, err := c.Read(buffer)
		if err != nil {
			c.Close()
			break
		}

		var b bytes.Buffer
		cmd, err := workers.NewCmd(string(bytes.Trim(buffer, "\x00")), &b)
		if err != nil {
			r.Error(err.Error())
			_, err = c.Write([]byte("error"))
		}

		err = cmd.Run()
		if err != nil {
			r.Error(err.Error())
			_, err = c.Write([]byte("error"))
		}

		_, err = c.Write(cmd.Stdout.Bytes())
		if err != nil {
			r.Error(err.Error())
		}
	}

	r.Info("Connection from %v closed.", c.RemoteAddr())
}

// RegisterAPIHandles to be used in http communication
func (r *Remote) RegisterAPIHandles() map[string]func(w http.ResponseWriter, r *http.Request) {
	cmds := make(map[string]func(w http.ResponseWriter, r *http.Request))

	cmds["start"] = r.cmdstart

	return cmds
}
