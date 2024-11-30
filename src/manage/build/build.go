package build

import (
	"context"
	"fmt"
	"github.com/charmbracelet/log"
	"io"
	"os"
	"os/exec"
	"path"
)

type Config struct {
	Ctx context.Context

	WorkingDir string
	OutputFile string

	Stdout io.Writer
	Stderr io.Writer
}

func GoPlugin(cnf *Config) error {
	return runGo(cnf, []string{
		"build",
		"-buildmode=plugin",
		fmt.Sprintf("-o=%s", cnf.OutputFile),
		".",
	}, []string{
		"CGO_ENABLED=1",
	}, "build plugin")
}

func GoFmt(cnf *Config) error {
	return runGo(cnf, []string{"fmt"}, []string{}, "fmt")
}

func runGo(cnf *Config, args []string, env []string, op string) error {
	binary, err := exec.LookPath("go")
	if err != nil {
		return err
	}

	if !path.IsAbs(cnf.WorkingDir) {
		return fmt.Errorf("working directory path should be absolute, got %s", cnf.WorkingDir)
	}

	environ := []string{
		fmt.Sprintf("HOME=%s", os.Getenv("HOME")),
		fmt.Sprintf("PATH=%s", os.Getenv("PATH")),
	}

	environ = append(environ, env...)

	cmd := exec.CommandContext(cnf.Ctx, binary, args...)
	cmd.Dir = cnf.WorkingDir
	cmd.Env = environ
	cmd.Stdout = cnf.Stdout
	cmd.Stderr = cnf.Stderr

	l := log.With("cmd", cmd.String(), "wd", cnf.WorkingDir, "action", op)
	l.Info("about to run")

	err = cmd.Run()

	if err != nil {
		l.With("err", err).Error("failed to run")
	} else {
		l.Info("cmd successful")
	}

	return err
}
