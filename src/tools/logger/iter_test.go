package logger

import (
	"bytes"
	"context"
	"errors"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/assert"
	"strings"
	"sync"
	"testing"
)

// Test that AllFmt logs all messages.
func TestIterationLogger_AllFmt(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	output := bytes.NewBufferString("")
	il := CreateIterationLogger(ctx, WithWriter(output))

	var wg sync.WaitGroup
	wg.Add(2)

	il.AllFmt()

	// Add a fork to ensure writing is done before asserting output
	il.fmt(func(msg Msg) bool {
		defer wg.Done()
		return false
	})

	il.Log(1, "test message 1")
	il.Log(2, "test message 2")

	// Stop the logger
	cancel()
	wg.Wait()

	// Validate output
	expected := strings.Join([]string{
		"[1] \x1b[0;34mINF\x1b[m: test message 1",
		"[2] \x1b[0;34mINF\x1b[m: test message 2",
		"",
	}, "\n")
	assert.Equal(t, expected, output.String())
}

// Test that FilterIterFmt only logs messages for a specific iteration.
func TestIterationLogger_FilterIterFmt(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	output := bytes.NewBufferString("")
	il := CreateIterationLogger(ctx, WithWriter(output))

	var wg sync.WaitGroup
	wg.Add(2)

	il.FilterIterFmt(1)

	// Add a fork to ensure writing is done before asserting output
	il.fmt(func(msg Msg) bool {
		defer wg.Done()
		return false
	})

	il.Log(1, "test message 1")
	il.Log(2, "test message 2")

	// Stop the logger
	cancel()
	wg.Wait()

	// Validate output
	expected := "[1] \x1b[0;34mINF\x1b[m: test message 1\n"
	assert.Equal(t, expected, output.String())
}

// Test dynamic fork addition and removal
func TestIterationLogger_RemoveFork(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	output := bytes.NewBufferString("")
	il := CreateIterationLogger(ctx, WithWriter(output))

	var wg sync.WaitGroup
	wg.Add(2)

	forkID := il.AllFmt()

	// Add a fork to ensure writing is done before asserting output
	il.fmt(func(msg Msg) bool {
		defer wg.Done()
		return false
	})

	il.Log(1, "test message 1")

	il.RemoveFork(forkID)

	il.Log(2, "test message 2")

	// Stop the logger
	cancel()
	wg.Wait()

	// Validate output
	expected := "[1] \x1b[0;34mINF\x1b[m: test message 1\n"
	assert.Equal(t, expected, output.String())
}

// Test LogPanic logs a recovered panic message.
func TestIterationLogger_LogPanic(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	output := bytes.NewBufferString("")
	il := CreateIterationLogger(ctx, WithWriter(output))

	var wg sync.WaitGroup
	wg.Add(1)

	il.AllFmt()

	// Add a fork to ensure writing is done before asserting output
	il.fmt(func(msg Msg) bool {
		defer wg.Done()
		return false
	})

	il.LogPanic(1, errors.New("simulated panic"))

	cancel()
	wg.Wait()

	expected := "[1] \x1b[0;31mPAN\x1b[m: simulated panic\n"
	assert.Equal(t, expected, output.String())
}

func TestIterationLogger_ContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	output := bytes.NewBufferString("")
	il := CreateIterationLogger(ctx, WithWriter(output))

	il.AllFmt()
	il.Log(1, "test before cancel")

	cancel()

	// Attempt to log after context cancellation
	il.Log(2, "test after cancel")

	expected := "[1] \x1b[0;34mINF\x1b[m: test before cancel\n"
	assert.Equal(t, expected, output.String())
}
