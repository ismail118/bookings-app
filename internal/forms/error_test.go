package forms

import "testing"

func TestErrors_Add(t *testing.T) {
	errs := errors{}
	errs.Add("e", "e")

	if len(errs) == 0 {
		t.Error("len should not 0")
	}
}

func TestErrors_Get(t *testing.T) {
	errs := errors{
		"e": []string{"e"},
	}
	e := errs.Get("e")

	if e == "" {
		t.Error("failed get")
	}

	errs = errors{}
	e = errs.Get("e")

	if e != "" {
		t.Error("failed get")
	}
}
