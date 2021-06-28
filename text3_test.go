package balaboba

import "testing"

func TestGet(t *testing.T) {
	c := New()

	// normal request
	gen, err := c.Get("123", NoStyle)
	if err != nil {
		t.Fatal(err)
	}
	if gen.Error != 0 {
		t.Log("bad error", gen.Error)
		t.FailNow()
	}
	if gen.BadQuery != 0 {
		t.Log("bad query", gen.BadQuery)
		t.FailNow()
	}

	// invalid style, but it's ok for api.
	gen, err = c.Get("123", Style(20))
	if err != nil {
		t.Fatal(err)
	}
	if gen.Error != 0 {
		t.Log("bad error", gen.Error)
		t.FailNow()
	}
	if gen.BadQuery != 0 {
		t.Log("bad query", gen.BadQuery)
		t.FailNow()
	}

	// bad query
	gen, err = c.Get(string([]rune{1093, 1091, 1081}), Style(20))
	if err != nil {
		t.Fatal(err)
	}
	if gen.Error != 0 {
		t.Log("bad error", gen.Error)
		t.FailNow()
	}
	if gen.BadQuery == 0 {
		t.Log("NO bad query")
	}
}

func TestOptions(t *testing.T) {
	c := New()

	err := c.Options()
	if err != nil {
		t.Fatal(err)
	}
}
