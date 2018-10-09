package recipe

import (
	"testing"
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

	// default time
	if !r.StartAt.IsZero() {
		t.Fail()
	}

	_, err = Build("../tests/etc/leprechaun/recipes/not_valid.yml")
	if err != nil {
		t.Error(err)
	}
}
