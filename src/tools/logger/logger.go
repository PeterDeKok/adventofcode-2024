package logger

import (
	"fmt"
	"github.com/charmbracelet/log"
	"os"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tools/env"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tools/exit"
)

type Logger = log.Logger

var logFp *os.File

// Init initializes a global file pointer to the log file specified in the environment.
// It is advisable for the enviroment variable to be an absolute filepath.
// The parent directory of the filepath should exist.
func Init() (*Logger, error) {
	e := env.MustGet()

	f, err := os.OpenFile(e.LogFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0660)
	if err != nil {
		return nil, err
	}
	logFp = f

	log.SetOutput(f)

	return log.Default(), nil
}

// MustInit returns a valid logger object, initialised with an output file.
// If the environment is invalid, it will output a help text to stdout and exit the program.
// If the file can't be opened, it will output an error to stdout and exit the program.
func MustInit() *Logger {
	if l, err := Init(); err == nil {
		return l
	} else {
		fmt.Printf("unable to open or create the log file: %v\n", err)

		panic(exit.ErrExitLogger)
	}
}
