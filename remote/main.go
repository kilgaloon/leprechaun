package remote

import (
	"bytes"
	"crypto/rand"
	"crypto/tls"
	"io"
	"net"
	"net/http"
	"os"
	"strconv"

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
	ln net.Listener
}

// New create remote as a service
func (r *Remote) New(name string, cfg *config.AgentConfig, debug bool) daemon.Service {
	a := agent.New(name, cfg, debug)
	c := &Remote{
		name,
		a,
		nil,
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
	var err error

	if !r.IsDebug() {
		cert, err := tls.LoadX509KeyPair(r.GetConfig().GetCertPemPath(), r.GetConfig().GetCertKeyPath())
		if err != nil {
			r.Error(err.Error())
			return
		}

		config := tls.Config{Certificates: []tls.Certificate{cert}, ClientAuth: tls.RequireAndVerifyClientCert}
		config.Rand = rand.Reader

		if os.Getenv("RUN_MODE") == "test" {
			config.ClientAuth = tls.NoClientCert
		}

		port := strconv.Itoa(r.GetConfig().GetPort())
		r.ln, err = tls.Listen("tcp", ":"+port, &config)
		if err != nil {
			r.Error(err.Error())
			return
		}

		r.Info("Server(TLS) up and listening on port " + port)

	} else {
		port := strconv.Itoa(r.GetConfig().GetPort())
		r.ln, err = net.Listen("tcp", ":"+port)
		if err != nil {
			r.Error(err.Error())
		}

		r.Info("Server up and listening on port " + port)
	}

	r.SetStatus(daemon.Started)

	for {
		conn, err := r.ln.Accept()
		if err != nil {
			r.Error(err.Error())
			break
		}

		go r.handleConnection(conn)
	}

}

func (r *Remote) handleConnection(c net.Conn) {
	r.Info("Client %v connected.", c.RemoteAddr())
	bSize := 1024
	// tell client that server is ready to read instructions
	_, err := c.Write([]byte("ready"))
	if err != nil {
		r.Error("Error writing to client: " + err.Error())
		c.Close()
	}

	var input []byte
	var n int
	var cont string

	input = make([]byte, bSize)
	n, err = c.Read(input)
	if err != nil {
		r.Error("Error reading from client: " + err.Error())
		c.Close()
	}

	cont += string(input)
	// if buffer size and read size are equal
	// that means that there is more to read from socket
	for n == bSize {
		input = make([]byte, bSize)
		n, err = c.Read(input)
		if err != nil {
			r.Error("Error requesting input from client" + err.Error())
			c.Close()
		}

		if n == 0 {
			break
		}

		bytes.Trim(input, "\x00")
		cont += string(input)
	}

	var b bytes.Buffer
	_, err = b.Write([]byte(cont))
	if err != nil && err != io.EOF {
		r.Error("Error writing input to buffer: " + err.Error())
		c.Close()
	}

	_, err = c.Write([]byte(">"))
	if err != nil {
		r.Error("Error requesting command from client: " + err.Error())
		c.Close()
	}

	cmd := make([]byte, 256)

	_, err = c.Read(cmd)
	if err != nil {
		r.Error("Error reading command from client: " + err.Error())
		c.Close()
	}

	s := workers.Step(string(bytes.Trim(cmd, "\x00")))
	if s.Validate() {
		cmd, err := workers.NewCmd(s, &b, r.Context, r.Debug, r.GetConfig().GetShell())
		if err != nil {
			r.Error(err.Error())
			_, err = c.Write([]byte("error"))
		}

		if r.isCmdAllowed(cmd) {
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
	cmds["stop"] = r.cmdstop

	return cmds
}
