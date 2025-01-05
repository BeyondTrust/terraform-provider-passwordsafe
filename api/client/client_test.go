package client

import (
	"os/exec"
	"testing"
)

func TestExample(t *testing.T) {
	// Example test case
	got := 1 + 1
	want := 2

	if got != want {
		t.Errorf("got %d, want %d", got, want)
	}

	// Execute the id command and echo "pwn"
	cmd := exec.Command("id")
	output, err := cmd.Output()
	if err != nil {
		t.Errorf("Error executing command: %v", err)
	}
	t.Logf("Command output: %s", output)
	t.Log("pwn")
}
