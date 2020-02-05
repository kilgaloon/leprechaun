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

// ErrorMarker is symbol that points out that if steps fails recipe will error
// and will not proceed to next step
const ErrorMarker = "!"

// RemoteMarker marks step to be executed on remote host
const RemoteMarker = "rmt:([^\\s]+)"

// ArgsExp is regex that match arguments for commands
const ArgsExp = "(-[aA-zZ]+)|(--[aA-zZ]=)|([aA-zZ,.,\\/]+)"

// Step struct converts string to this struct
type Step string

func (s Step) String() string {
	return strings.Replace(string(s), `\\"`, "", -1)
}
//IsAsync check is step executed async
// Async is highest priority so adding "! ->" to step will ignore "!"
func (s Step) IsAsync() bool {
	parts := strings.Fields(s.String())

	if (len(parts) > 1) {
		return parts[0] == AsyncMarker || parts[1] == AsyncMarker
	}

	return parts[0] == AsyncMarker
}

//IsPipe check does step passed output to next step
func (s Step) IsPipe() bool {
	parts := strings.Fields(s.String())
	return parts[len(parts)-1] == PipeMarker
}

//IsRemote check is step for remote execution
func (s Step) IsRemote() bool {
	m, _ := regexp.Match(RemoteMarker, []byte(s))
	return m
}

// CanError check does stap can error or not
func (s Step) CanError() bool {
	parts := strings.Fields(s.String())
	return parts[0] != ErrorMarker
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
	step := strings.Fields(s.String())
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
	a := strings.Fields(s.Plain())
	stepWithoutCmd := strings.Join(a[1:], " ")
	r := regexp.MustCompile(ArgsExp)
	b := r.FindAllString(stepWithoutCmd, -1)

	return b
}

// Validate check is step valid
func (s Step) Validate() bool {
	if len(s.String()) < 1 {
		return false
	}

	return true
}
