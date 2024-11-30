package assert

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"testing"
)

func CatchExit(t *testing.T, code int, fn func()) {
	t.Helper()

	if c := os.Getenv("AOC_TEST_EXIT_NON_ZERO"); c != "" && c != fmt.Sprintf("%d", code) {
		return
	} else if c == fmt.Sprintf("%d", code) {
		fn()
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run="+t.Name())
	cmd.Env = append(os.Environ(), fmt.Sprintf("AOC_TEST_EXIT_NON_ZERO=%d", code))
	err := cmd.Run()
	var e *exec.ExitError
	if errors.As(err, &e) && !e.Success() && e.ExitCode() == code {
		return
	}
	t.Fatalf("process ran with err %v, want exit status %d", err, code)
}
