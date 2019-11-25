package remote

import (
	"bytes"
	"crypto/rand"
	"crypto/tls"
	"io"
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

func (r Remote) isCmdAllowed(cmd *workers.Cmd) bool {
	cmds := r.GetConfig().Cfg.Section("").Key(r.GetName() + ".allowed_commands").Strings(",")

	for _, c := range cmds {
		if c == cmd.Step.Name() {
			return true
		}
	}

	return false
}

// Start remote service
func (r *Remote) Start() {
	var ln net.Listener
	var err error

	if !r.IsDebug() {
		cert, err := tls.LoadX509KeyPair(r.GetConfig().GetCertPemPath(), r.GetConfig().GetCertKeyPath())
		if err != nil {
			r.Error(err.Error())
		}

		config := tls.Config{Certificates: []tls.Certificate{cert}, ClientAuth: tls.NoClientCert}
		config.Rand = rand.Reader

		ln, err = tls.Listen("tcp", ":11402", &config)
		if err != nil {
			r.Error(err.Error())
		}

		r.Info("Server(TLS) up and listening on port 11402")

	} else {
		ln, err = net.Listen("tcp", ":11402")
		if err != nil {
			r.Error(err.Error())
		}
	}

	r.SetStatus(daemon.Started)

	for {
		conn, err := ln.Accept()
		if err != nil {
			r.Error(err.Error())
			continue
		}

		go r.handleConnection(conn)
	}

}

func (r *Remote) handleConnection(c net.Conn) {
	r.Info("Client %v connected.", c.RemoteAddr())
	bSize := 1024
	// tell client that server is ready to read instructions
	c.Write([]byte("ready"))

	var input []byte
	var n int
	var cont string
	var err error

	input = make([]byte, bSize)
	n, err = c.Read(input)
	cont += string(input)
	// if buffer size and read size are equal
	// that means that there is more to read from socket
	for n == bSize {
		input = make([]byte, bSize)
		n, err = c.Read(input)
		if err != nil {
			c.Close()
		}

		if n == 0 {
			break
		}

		cont += string(input)
	}

	var b bytes.Buffer
	_, err = b.Write(bytes.Trim(input, "\x00"))
	if err != nil && err != io.EOF {
		c.Close()
	}

	_, err = c.Write([]byte(">"))
	cmd := make([]byte, 256)

	_, err = c.Read(cmd)
	if err != nil {
		c.Close()
	}

	s := workers.Step(string(bytes.Trim(cmd, "\x00")))
	if s.Validate() {
		cmd, err := workers.NewCmd(s, &b, r.Context, r.Debug)
		if r.isCmdAllowed(cmd) {
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
		} else {
			r.Error("Command not allowed %s", s.Name())
			_, err = c.Write([]byte("error"))
		}
	}

	r.Info("Connection from %v closed.", c.RemoteAddr())
	c.Close()
}

// RegisterAPIHandles to be used in http communication
func (r *Remote) RegisterAPIHandles() map[string]func(w http.ResponseWriter, r *http.Request) {
	cmds := make(map[string]func(w http.ResponseWriter, r *http.Request))

	cmds["start"] = r.cmdstart

	return cmds
}
