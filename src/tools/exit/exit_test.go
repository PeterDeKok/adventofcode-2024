package exit

import (
	"errors"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/assert"
	"testing"
)

func TestErrExit_Error(t *testing.T) {
	assert.TypeOf[error](t, ErrExitEnv)
	assert.Equal(t, "exit 11: env", ErrExitEnv.Error())
	assert.TypeOf[error](t, ErrExitLogger)
	assert.Equal(t, "exit 12: logger", ErrExitLogger.Error())
	assert.TypeOf[error](t, ErrExitManager)
	assert.Equal(t, "exit 13: manager", ErrExitManager.Error())
}

func TestPanicToExit(t *testing.T) {
	assert.CatchExit(t, ErrExitEnv.Code, func() {
		defer PanicToExit()
		panic(ErrExitEnv)
	})

	assert.CatchExit(t, ErrExitLogger.Code, func() {
		defer PanicToExit()
		panic(ErrExitLogger)
	})

	assert.CatchExit(t, ErrExitManager.Code, func() {
		defer PanicToExit()
		panic(ErrExitManager)
	})

	assert.CatchExit(t, 255, func() {
		defer PanicToExit()
		panic(errors.New("not a defined error"))
	})
}
