package workers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultStep(t *testing.T) {
	s := Step("echo \"default step\"")

	if !s.CanError() {
		t.Fatal("Command can error but evaluated to cant")
	}

	if s.IsAsync() {
		t.Fail()
	}

	if !s.Validate() {
		t.Fail()
	}

	if s.Plain() != "echo \"default step\"" {
		t.Fail()
	}
}

func TestAsyncStep(t *testing.T) {
	s := Step("-> echo \"default step\"")

	if !s.CanError() {
		t.Fatal("Command can error but evaluated to cant")
	}

	if !s.IsAsync() {
		t.Fail()
	}

	if s.IsPipe() {
		t.Fail()
	}

	if !s.Validate() {
		t.Fail()
	}

	if s.Plain() != "echo \"default step\"" {
		t.Fail()
	}
}

func TestPipeStep(t *testing.T) {
	s := Step("echo \"default step\" }>")

	if !s.CanError() {
		t.Fatal("Command can error but evaluated to cant")
	}

	if s.IsAsync() {
		t.Fail()
	}

	if !s.IsPipe() {
		t.Fail()
	}

	if !s.Validate() {
		t.Fail()
	}

	if s.Plain() != "echo \"default step\"" {
		t.Fail()
	}
}

func TestAsyncPipeStep(t *testing.T) {
	s := Step("-> echo \"default step\" }>")

	if !s.CanError() {
		t.Fatal("Command can error but evaluated to cant")
	}

	if !s.IsAsync() {
		t.Fail()
	}

	if !s.IsPipe() {
		t.Fail()
	}

	if !s.Validate() {
		t.Fail()
	}

	if s.Plain() != "echo \"default step\"" {
		t.Fail()
	}
}

func TestRemoteStep(t *testing.T) {
	s := Step("rmt:host echo \"default step\"")

	if !s.CanError() {
		t.Fatal("Command can error but evaluated to cant")
	}

	if s.IsAsync() {
		t.Fail()
	}

	if s.IsPipe() {
		t.Fail()
	}

	if !s.Validate() {
		t.Fail()
	}

	if !s.IsRemote() {
		t.Fail()
	}

	if s.RemoteHost() != "host" {
		t.Fail()
	}

	if s.Plain() != "echo \"default step\"" {
		t.Fail()
	}
}

func TestRemoteAsyncStep(t *testing.T) {
	s := Step("-> rmt:host echo \"default step\"")

	if !s.CanError() {
		t.Fatal("Command can error but evaluated to cant")
	}

	if !s.IsAsync() {
		t.Fail()
	}

	if s.IsPipe() {
		t.Fail()
	}

	if !s.Validate() {
		t.Fail()
	}

	if !s.IsRemote() {
		t.Fail()
	}

	if s.RemoteHost() != "host" {
		t.Fail()
	}

	if s.Plain() != "echo \"default step\"" {
		t.Fail()
	}
}

func TestRemotePipeStep(t *testing.T) {
	s := Step("rmt:host echo \"default step\" }>")

	if !s.CanError() {
		t.Fatal("Command can error but evaluated to cant")
	}

	if s.IsAsync() {
		t.Fail()
	}

	if !s.IsPipe() {
		t.Fail()
	}

	if !s.Validate() {
		t.Fail()
	}

	if !s.IsRemote() {
		t.Fail()
	}

	if s.RemoteHost() != "host" {
		t.Fail()
	}

	if s.Plain() != "echo \"default step\"" {
		t.Fail()
	}
}

func TestRemoteAsyncPipeStep(t *testing.T) {
	s := Step("-> rmt:host echo \"default step\" }>")

	if !s.CanError() {
		t.Fatal("Command can error but evaluated to cant")
	}

	if !s.IsAsync() {
		t.Fail()
	}

	if !s.IsPipe() {
		t.Fail()
	}

	if !s.Validate() {
		t.Fail()
	}

	if !s.IsRemote() {
		t.Fail()
	}

	if s.RemoteHost() != "host" {
		t.Fail()
	}

	if s.Plain() != "echo \"default step\"" {
		t.Fail()
	}
}

func TestName(t *testing.T) {
	s := Step("echo \"default step\"")

	if s.Name() != "echo" {
		t.Fail()
	}

	s = Step("/usr/share/echo \"default step\"")
	if s.Name() != "echo" {
		t.Fail()
	}

	if s.FullName() != "/usr/share/echo" {
		t.Fail()
	}
}

func TestArgs(t *testing.T) {
	s := Step("/usr/share/echo arg --output=config.conf --ini=\"/etc/prog/def.ini\" -flag")
	if len(s.Args()) < 1 {
		t.Fail()
	}

	assert.Equal(t, s.Args()[0], "arg")
	assert.Equal(t, "-output", s.Args()[1])
	assert.Equal(t, "config.conf", s.Args()[2])
	assert.Equal(t, "-ini", s.Args()[3])
	assert.Equal(t, "/etc/prog/def.ini", s.Args()[4])
	assert.Equal(t, "-flag", s.Args()[5])
}

func TestCantError(t *testing.T) {
	s := Step("! -> rmt:host echo \"default step\" }>")

	if s.CanError() {
		t.Fatal("Command cant error but evaluated to can")
	}

	if !s.IsAsync() {
		t.Fatal("Command is async but evaluated as not")
	}

	if !s.IsPipe() {
		t.Fatal("Command is piped but evaluated as not")
	}

	if !s.Validate() {
		t.Fatal("Command is valid but evaluated as not")
	}

	if !s.IsRemote() {
		t.Fatal("Command is remote but evaluated as not")
	}

	if s.RemoteHost() != "host" {
		t.Fatal("Command is remote but evaluated host is not valid")
	}

	if s.Plain() != "echo \"default step\"" {
		t.Fatal("Plain command not correct")
	}
}

func TestCanError(t *testing.T) {
	s := Step("-> ! rmt:host echo \"default step\" }>")

	if !s.CanError() {
		t.Fatal("Command cant error but evaluated to can")
	}

	if !s.IsAsync() {
		t.Fatal("Command is async but evaluated as not")
	}

	if !s.IsPipe() {
		t.Fatal("Command is piped but evaluated as not")
	}

	if !s.Validate() {
		t.Fatal("Command is valid but evaluated as not")
	}

	if !s.IsRemote() {
		t.Fatal("Command is remote but evaluated as not")
	}

	if s.RemoteHost() != "host" {
		t.Fatal("Command is remote but evaluated host is not valid")
	}

	if s.Plain() != "echo \"default step\"" {
		t.Fail()
	}
}
