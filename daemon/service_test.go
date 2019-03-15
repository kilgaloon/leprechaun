package daemon

import "testing"

func TestServiceStatus(t *testing.T) {
	started := ServiceStatus(1)
	stopped := ServiceStatus(2)
	paused := ServiceStatus(3)
	unknown := ServiceStatus(10)

	if started.String() != "Started" {
		t.Fail()
	}

	if stopped.String() != "Stopped" {
		t.Fail()
	}

	if paused.String() != "Paused" {
		t.Fail()
	}

	if unknown.String() != "Unknown" {
		t.Fail()
	}
}
