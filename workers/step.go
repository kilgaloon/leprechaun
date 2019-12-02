package workers

import (
	"os"
	"regexp"
	"strings"
)

// AsyncMarker is string in step that we use to know
// does that command need to be done async
const AsyncMarker = "->"

// PipeMarker is string that we mark when we want to pipe output of results
// from step command to another step command
const PipeMarker = "}>"

// RemoteMarker marks step to be executed on remote host
const RemoteMarker = "rmt:([^\\s]+)"

// ArgsExp is regex that match arguments for commands
const ArgsExp = "(-[aA-zZ]+)|(--[aA-zZ]+-[aA-zZ]+=[\"]?[aA-zZ,.,]+[\"]?)|(\".*\")|([aA-zZ,.,]+)"

// Step struct converts string to this struct
type Step string

//IsAsync check is step executed async
func (s Step) IsAsync() bool {
	parts := strings.Fields(string(s))
	return parts[0] == AsyncMarker
}

//IsPipe check does step passed output to next step
func (s Step) IsPipe() bool {
	parts := strings.Fields(string(s))
	return parts[len(parts)-1] == PipeMarker
}

//IsRemote check is step for remote execution
func (s Step) IsRemote() bool {
	m, _ := regexp.Match(RemoteMarker, []byte(s))
	return m
}

// RemoteHost extracts name of host from step "rmt:host"
func (s Step) RemoteHost() string {
	r := regexp.MustCompile(RemoteMarker)
	rmt := r.Find([]byte(s))
	host := strings.Split(string(rmt), ":")

	return host[1]
}

// Plain return command without markers and our syntax
// ex: -> echo "Test" }> will result in echo "Test"
func (s Step) Plain() string {
	step := strings.Fields(string(s))
	var a, b int

	if step[0] == AsyncMarker {
		a = 1
	}

	b = len(step)
	if step[b-1] == PipeMarker {
		b = len(step) - 1
	}

	stepTrimmed := strings.Join(step[a:b], " ")

	r := regexp.MustCompile(RemoteMarker)
	step = r.Split(stepTrimmed, 2)

	stepTrimmed = strings.TrimLeft(step[len(step)-1], " ")

	return stepTrimmed
}

// TODO: This returns name of command not name of step, change name to something appropiate

// Name returns name of command
func (s Step) Name() string {
	a := strings.Fields(s.Plain())
	b := strings.Split(a[0], string(os.PathSeparator))

	return b[len(b)-1]
}

// TODO: This returns name of command not name of step, change name to something appropiate

// FullName to to command included with Path seperators
func (s Step) FullName() string {
	a := strings.Fields(s.Plain())

	return a[0]
}

// Args extract arguments for command
func (s Step) Args() []string {
	r := regexp.MustCompile(ArgsExp)
	b := r.FindAllString(s.Plain(), -1)

	return b[1:]
}

// Validate check is step valid
func (s Step) Validate() bool {
	if len(string(s)) < 1 {
		return false
	}

	return true
}
