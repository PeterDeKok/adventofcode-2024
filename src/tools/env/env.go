package env

import (
	"errors"
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"os"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tools/exit"
	"strings"
	"time"
)

type Env struct {
	LogFile              string
	PuzzlesDir           string
	SessionCookieValue   string
	SessionCookieExpires time.Time
}

type envResult[T any] struct {
	key         string
	description string
	v           T
	err         error
	errExtra    error
	ok          bool
}

type envInternal struct {
	ok bool

	LogFile              envResult[string]
	PuzzlesDir           envResult[string]
	SessionCookieValue   envResult[string]
	SessionCookieExpires envResult[time.Time]
}

var env *Env

var (
	ErrDirInaccessible  = errors.New("directory inaccessible")
	ErrFileInaccessible = errors.New("file inaccessible")
	ErrFileIsDir        = errors.New("file is a directory")
	ErrDirIsFile        = errors.New("directory is a file")
	ErrMissing          = errors.New("missing")
	ErrInvalid          = errors.New("invalid environment")
	ErrInvalidTime      = errors.New("invalid time format")
)

var (
	StyleErr = lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(ansi.Red))
)

// Get retrieves and validates the relevant environment.
// It will return an error and help text in case of failure.
func Get() (e *Env, help string, err error) {
	if env != nil {
		return env, "", nil
	}

	ie := &envInternal{
		ok: true,

		LogFile:              envResult[string]{key: "AOC_LOG_FILE", description: "the filepath to the log file. it is advisible to use an absolute path, however it will load relative paths as well. the parent directory should exist"},
		PuzzlesDir:           envResult[string]{key: "AOC_PUZZLES_DIR", description: "the filepath to the puzzles directory. it is advisible to use an absolute path, however it will load relative paths as well. the directory should exist"},
		SessionCookieExpires: envResult[time.Time]{key: "AOC_SESSION_COOKIE_EXPIRES", description: "the expiration date of the aoc session cookie"},
		SessionCookieValue:   envResult[string]{key: "AOC_SESSION_COOKIE_VALUE", description: "the content of the aoc session cookie"},
	}

	{
		v := strings.TrimSpace(os.Getenv("AOC_LOG_FILE"))
		if len(v) == 0 {
			ie.ok = false
			ie.LogFile.err = ErrMissing
		} else if fi, err := os.Stat(v); err != nil && !os.IsNotExist(err) {
			ie.ok = false
			ie.LogFile.err = ErrFileInaccessible
			ie.LogFile.errExtra = err
		} else if err == nil && fi.IsDir() {
			ie.ok = false
			ie.LogFile.err = ErrFileInaccessible
			ie.LogFile.errExtra = ErrFileIsDir
		} else {
			ie.LogFile.v = v
			ie.LogFile.ok = true
		}
	}

	{
		v := strings.TrimSpace(os.Getenv("AOC_PUZZLES_DIR"))
		if len(v) == 0 {
			ie.ok = false
			ie.PuzzlesDir.err = ErrMissing
		} else if fi, err := os.Stat(v); err != nil && !os.IsNotExist(err) {
			ie.ok = false
			ie.PuzzlesDir.err = ErrDirInaccessible
			ie.PuzzlesDir.errExtra = err
		} else if err == nil && !fi.IsDir() {
			ie.ok = false
			ie.PuzzlesDir.err = ErrDirInaccessible
			ie.PuzzlesDir.errExtra = ErrDirIsFile
		} else {
			ie.PuzzlesDir.v = v
			ie.PuzzlesDir.ok = true
		}
	}

	{
		vTmp := strings.TrimSpace(os.Getenv("AOC_SESSION_COOKIE_EXPIRES"))
		if len(vTmp) == 0 {
			ie.ok = false
			ie.SessionCookieExpires.err = ErrMissing
		} else if v, err := time.Parse(time.RFC1123, vTmp); err != nil {
			ie.ok = false
			ie.SessionCookieExpires.err = ErrInvalidTime
			ie.SessionCookieExpires.errExtra = err
		} else {
			ie.SessionCookieExpires.v = v
			ie.SessionCookieExpires.ok = true
		}
	}

	{
		v := strings.TrimSpace(os.Getenv("AOC_SESSION_COOKIE_VALUE"))
		if len(v) == 0 {
			ie.ok = false
			ie.SessionCookieValue.err = ErrMissing
		} else {
			ie.SessionCookieValue.v = v
			ie.SessionCookieValue.ok = true
		}
	}

	e = &Env{
		LogFile:              ie.LogFile.v,
		PuzzlesDir:           ie.PuzzlesDir.v,
		SessionCookieValue:   ie.SessionCookieValue.v,
		SessionCookieExpires: ie.SessionCookieExpires.v,
	}

	if !ie.ok {
		return e, fmt.Sprintf("environment\n\n%s%s%s%s",
			ie.LogFile,
			ie.PuzzlesDir,
			ie.SessionCookieExpires,
			ie.SessionCookieValue,
		), ErrInvalid
	}

	return e, "", nil
}

// MustGet returns a valid environment object.
// If the environment is invalid, it will output a help text to stdout
// and exit the program.
func MustGet() *Env {
	if env, help, err := Get(); err == nil {
		return env
	} else if len(help) > 0 {
		fmt.Println(help)

		panic(exit.ErrExitEnv)
	} else {
		panic(err)
	}
}

func (s envResult[T]) String() string {
	var errStr string

	if s.err != nil {
		errStr = "err: " + s.err.Error()

		if s.errExtra != nil {
			errStr += ": " + s.errExtra.Error()
		}

		errStr = ansi.Hardwrap(errStr, 80-4, true)
		errStr = StyleErr.Render("    "+strings.ReplaceAll(errStr, "\n", "\n    ")) + "\n"
	}

	descStr := ansi.Hardwrap(s.description, 80-4, true)
	descStr = strings.ReplaceAll(descStr, "\n", "\n    ")

	return fmt.Sprintf("%s\n    %s\n%s\n", s.key, descStr, errStr)
}
