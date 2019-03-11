package recipe

import (
	"testing"
	"time"
)

func TestBuild(t *testing.T) {
	_, err := Build("../tests/etc/leprechaun/recipes/schedule.yml")
	if err != nil {
		t.Error(err)
	}

	_, err = Build("../tests/etc/leprechaun/recipes/not_exists.yml")
	if err == nil {
		t.Error(err)
	}

	r, err := Build("../tests/etc/leprechaun/recipes/hook.yml")
	if err != nil {
		t.Error(err)
	}

	if r.GetName() != r.Name {
		t.Fail()
	}

	l1 := len(r.GetSteps())
	l2 := len(r.Steps)
	if l1 != l2 {
		t.Fatal("Steps on same")
	}

	// default time
	if !r.GetStartAt().IsZero() {
		t.Fail()
	}

	time := time.Now()
	r.SetStartAt(time)
	if !time.Equal(r.GetStartAt()) {
		t.Fatal("Time not equal")
	}

	_, err = Build("../tests/etc/leprechaun/recipes/not_valid.yml")
	if err != nil {
		t.Error(err)
	}
}
