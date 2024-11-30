package logger

import (
	"os"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/assert"
	"testing"
)

func TestInit_NoEnv(t *testing.T) {
	assert.NoErr(t, os.Unsetenv("AOC_LOG_FILE"))

	assert.ShouldPanic(t, func() {
		_, _ = Init()
	})
}

func TestInit(t *testing.T) {
	assert.NoErr(t, os.Setenv("AOC_LOG_FILE", "/tmp/aoc-test-logger.log.txt"))

	t.Cleanup(func() {
		assert.NoErr(t, os.Remove("/tmp/aoc-test-logger.log.txt"))
	})

	l, err := Init()

	assert.NoErr(t, err)
	assert.NotNil(t, l)
	assert.TypeOf[*Logger](t, l)

	l.Error("tools.logger.TestInit")

	f, err := os.ReadFile("/tmp/aoc-test-logger.log.txt")

	assert.NoErr(t, err)
	assert.StrCount(t, "\n", 1, string(f))
	assert.EndsWith(t, "ERRO tools.logger.TestInit\n", string(f))
}
