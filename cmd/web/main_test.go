package main

import (
	"os"
	"testing"
)

func TestRun(t *testing.T) {
	os.Args = []string{"test", "-dbname", "bookings_app", "-dbuser", "postgres", "-dbpass", "postgres", "-cache", "false", "-production", "false"}
	_, err := run()
	if err != nil {
		t.Error("failed run()")
	}
}
