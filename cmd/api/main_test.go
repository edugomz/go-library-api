package main

import (
	"os/exec"
	"strings"
	"testing"
)

// The server must never run migrations on a normal boot (Cloud Run can start
// several instances concurrently, and concurrent migration runs race).
// Migrations are only run via the explicit --migrate-only flag, so the flag
// must stay documented and wired into the entrypoint.
func TestMigrateOnlyFlagRegistered(t *testing.T) {
	out, _ := exec.Command("go", "run", ".", "--help").CombinedOutput()

	if !strings.Contains(string(out), "migrate-only") {
		t.Fatalf("expected --migrate-only flag to be documented in --help output, got:\n%s", out)
	}
}
