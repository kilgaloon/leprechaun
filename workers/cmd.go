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

type cmd struct {
	stdin  bytes.Buffer
	cmd    *exec.Cmd
	stdout bytes.Buffer
	pipe   bool
}

func (c *cmd) Run() error {
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
	c.cmd.Stdout = &c.stdout
	c.cmd.Stderr = &stderr

	err := c.cmd.Run()

	return err
}

// NewCmd build new command and prepare it to be run
func newCmd(step string, i bytes.Buffer) (*cmd, error) {
	cmd := &cmd{
		stdin: i,
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
