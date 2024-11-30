package exit

import (
	"fmt"
	"os"
)

type errExit struct {
	Code int
	Desc string
}

var (
	ErrExitEnv     = errExit{11, "env"}
	ErrExitLogger  = errExit{12, "logger"}
	ErrExitManager = errExit{13, "manager"}
)

func (err errExit) Error() string {
	return fmt.Sprintf("exit %d: %s", err.Code, err.Desc)
}

// PanicToExit allows a critical error to exit the program with a custom exit code,
// while still allowing other deferred cleanup to take place, like terminal background
// reset.
// This should be used instead of os.Exit.
func PanicToExit() {
	if r := recover(); r != nil {
		switch r := r.(type) {
		case errExit:
			fmt.Println(r)
			os.Exit(r.Code)
		default:
			fmt.Println("unknown panic: ", r)
			os.Exit(255)
		}
	}
}
