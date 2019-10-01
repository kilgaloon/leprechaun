package workers

import (
	"bytes"
	"io"
	"os/exec"
	"strings"
)

// PipeMarker is string that we mark when we want to pipe output of results
// from step command to another step command
const PipeMarker = "}>"

// Cmd represents command that can be run
type Cmd struct {
	stdin  bytes.Buffer
	cmd    *exec.Cmd
	Stdout bytes.Buffer
	pipe   bool
}

// Run command and returns errors if any
func (c *Cmd) Run() error {
	if &c.stdin != nil {
		in, err := c.cmd.StdinPipe()
		if err != nil {
			return err
		}

		_, err = io.WriteString(in, string(c.stdin.Bytes()))
		if err != nil {
			return err
		}
	}

	var stderr bytes.Buffer
	c.cmd.Stdout = &c.Stdout
	c.cmd.Stderr = &stderr

	err := c.cmd.Run()

	return err
}

// NewCmd build new command and prepare it to be run
func NewCmd(step string, i *bytes.Buffer) (*Cmd, error) {
	cmd := &Cmd{
		stdin: *i,
		pipe:  false,
	}

	s := strings.Fields(step)
	if s[len(s)-1] == PipeMarker {
		cmd.pipe = true
		step = strings.Join(s[:(len(s)-1)], " ")
	}

	cmd.cmd = exec.Command("bash", "-c", step)

	return cmd, nil
}
