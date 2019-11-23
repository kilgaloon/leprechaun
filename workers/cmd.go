package workers

import (
	"bytes"
	"crypto/rand"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net"
	"os/exec"

	"github.com/kilgaloon/leprechaun/context"
)

// Cmd represents command that can be run
type Cmd struct {
	ctx    *context.Context
	Stdin  bytes.Buffer
	Cmd    *exec.Cmd
	Stdout bytes.Buffer
	Step   Step
	Remote bool
	pipe   bool
	Debug  bool
}

// Run command and returns errors if any
func (c *Cmd) Run() error {
	if c.Step.IsRemote() {
		return c.runRemote()
	}

	if &c.Stdin != nil {
		in, err := c.Cmd.StdinPipe()
		if err != nil {
			return err
		}

		w := string(bytes.Trim(c.Stdin.Bytes(), "\x00"))
		_, err = io.WriteString(in, w)
		if err != nil {
			return err
		}
	}

	var stderr bytes.Buffer
	c.Cmd.Stdout = &c.Stdout
	c.Cmd.Stderr = &stderr

	err := c.Cmd.Run()

	return err
}

func (c *Cmd) runRemote() (err error) {
	host := net.JoinHostPort(c.Step.RemoteHost(), "11402")
	var conn net.Conn

	if !c.Debug {
		cert, err := tls.LoadX509KeyPair(
			c.ctx.GetVar("pem_file_path").GetValue(),
			c.ctx.GetVar("key_file_path").GetValue(),
		)

		if err != nil {
			return err
		}

		config := tls.Config{Certificates: []tls.Certificate{cert}}
		config.Rand = rand.Reader

		conn, err = tls.Dial("tcp", host, &config)
		if err != nil {
			return err
		}
	} else {
		conn, err = net.Dial("tcp", host)
		if err != nil {
			return err
		}
	}

	message := make([]byte, 5)
	// listen for message
	_, err = conn.Read(message)
	msg := string(message)
	if msg != "ready" {
		err = errors.New("Failed to get back from server")
	}

	m := c.Stdin.Bytes()
	var n int
	for n < 1 {
		if m == nil {
			m = []byte("\n")
		}

		n, err = conn.Write(m)
		if err != nil {
			return err
		}
	}

	// expecting server to respond with ">" which means
	// that is waiting for command
	message = make([]byte, 1)
	// listen for message
	_, err = conn.Read(message)
	msg = string(message)
	if msg != ">" {
		err = errors.New("Failed to get back from server")
	}

	_, err = conn.Write([]byte(c.Step.Plain()))
	if err != nil {
		return
	}

	message = make([]byte, 1024)
	// listen for message
	var cont string
	n, err = conn.Read(message)
	cont += string(message)
	for n == 1024 {
		message = make([]byte, 1024)
		n, err = conn.Read(message)
		if err != nil {
			return
		}

		cont += string(message)
	}

	_, err = c.Stdout.Write(message)

	fmt.Printf("Server returned output: %s \n", message)

	return
}

// NewCmd build new command and prepare it to be run
func NewCmd(step Step, i *bytes.Buffer, ctx *context.Context, debug bool) (*Cmd, error) {
	cmd := &Cmd{
		ctx:   ctx,
		Stdin: *i,
		Step:  step,
		pipe:  false,
		Debug: debug,
	}

	if cmd.Step.IsPipe() {
		cmd.pipe = true
	}

	cmd.Cmd = exec.Command("bash", "-c", cmd.Step.Plain())

	return cmd, nil
}

// NewRemoteCmd creates new command to be executed remotely
func NewRemoteCmd(step Step, i *bytes.Buffer, ctx *context.Context, debug bool) (*Cmd, error) {
	cmd := &Cmd{
		ctx:   ctx,
		Stdin: *i,
		Step:  step,
		pipe:  false,
		Debug: debug,
	}

	if cmd.Step.IsPipe() {
		cmd.pipe = true
	}

	return cmd, nil
}
