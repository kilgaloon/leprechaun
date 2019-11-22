package workers

import (
	"testing"
)

func TestDefaultStep(t *testing.T) {
	s := Step("echo \"default step\"")

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
}
