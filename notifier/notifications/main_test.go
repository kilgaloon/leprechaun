package notifications

import "testing"

func TestOptions(t *testing.T) {
	o := Options{
		Title: "title",
		Body:  "body",
	}

	if o.GetBody() != "body" {
		t.Fail()
	}

	if o.GetTitle() != "title" {
		t.Fail()
	}

	o1 := Options{
		Title: "",
		Body:  "",
	}

	if o1.GetTitle() != "Leprechaun notification" {
		t.Fail()
	}
}
