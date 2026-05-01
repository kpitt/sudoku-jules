package puzzle

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestErrPuzzleFormat(t *testing.T) {
	err := errPuzzleFormat("test %d", 1)
	expected := "invalid puzzle format: test 1"
	if err.Error() != expected {
		t.Errorf("expected %q, got %q", expected, err.Error())
	}
}

func TestErrPuzzleState(t *testing.T) {
	err := errPuzzleState("test %d", 2)
	expected := "invalid puzzle state: test 2"
	if err.Error() != expected {
		t.Errorf("expected %q, got %q", expected, err.Error())
	}
}

func TestFatalError(t *testing.T) {
	if os.Getenv("BE_CRASHER") == "1" {
		fatalError("test", "message")
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestFatalError")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
	out, err := cmd.CombinedOutput()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		if !strings.Contains(string(out), "error: test: message") {
			t.Errorf("expected output to contain %q, got %q", "error: test: message", string(out))
		}
		return
	}
	t.Fatalf("process ran with err %v, want exit status 1", err)
}

func TestPuzzleStateError(t *testing.T) {
	if os.Getenv("BE_CRASHER") == "1" {
		puzzleStateError("test message")
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestPuzzleStateError")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
	out, err := cmd.CombinedOutput()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		if !strings.Contains(string(out), "error: invalid puzzle state: test message") {
			t.Errorf("expected output to contain %q, got %q", "error: invalid puzzle state: test message", string(out))
		}
		return
	}
	t.Fatalf("process ran with err %v, want exit status 1", err)
}
