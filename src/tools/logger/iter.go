package logger

import (
	"context"
	"fmt"
	"io"
	"os"
	"sync/atomic"
)

type Msg struct {
	Iter   int
	Level  Level
	Format string
	Args   []interface{}
}

type IterationLogger struct {
	ctx   context.Context
	ch    chan Msg
	forkI uint32
	fork  map[uint32]func(msg Msg)
	w     io.Writer
}

func CreateIterationLogger(ctx context.Context, opts ...Option) *IterationLogger {
	il := &IterationLogger{
		ctx:  ctx,
		ch:   make(chan Msg),
		fork: make(map[uint32]func(msg Msg)),
		w:    os.Stdout,
	}

	for _, opt := range opts {
		opt(il)
	}

	go il.listen()

	return il
}

func CreateIterationLoggerWithWriter(ctx context.Context, w io.Writer) *IterationLogger {
	il := &IterationLogger{
		ctx:  ctx,
		ch:   make(chan Msg),
		fork: make(map[uint32]func(msg Msg)),
		w:    w,
	}

	go il.listen()

	return il
}

// listen retrieves a message from the channel and calls all forks.
// A fork is a function which can process a log message.
// listen is intended to be ran in a goroutine.
func (il *IterationLogger) listen() {
	if il == nil {
		panic("logger borked")
	}

	ctx := il.ctx
	if ctx == nil {
		panic("logger ctx borked")
	}

	for {
		select {
		case <-ctx.Done():
			close(il.ch)
			return
		case msg, ok := <-il.ch:
			if !ok {
				return
			}

			for _, fn := range il.fork {
				fn(msg)
			}
		}
	}
}

// fmt prints messages to stdout when a filter has determined a message should be printed.
// When the filter is nil, it will print all messages.
func (il *IterationLogger) fmt(filter func(msg Msg) bool) uint32 {
	i := atomic.AddUint32(&il.forkI, 1)

	if filter == nil {
		il.fork[i] = func(msg Msg) {
			_, err := fmt.Fprintf(il.w, "[%d] %s: %s\n", msg.Iter, msg.Level.Fmt(), fmt.Sprintf(msg.Format, msg.Args...))
			if err != nil {
				panic(fmt.Errorf("failed to write log message: [%d] %s: %s\n", msg.Iter, msg.Level.Fmt(), fmt.Sprintf(msg.Format, msg.Args...)))
			}
		}
	} else {
		il.fork[i] = func(msg Msg) {
			if !filter(msg) {
				return
			}

			_, err := fmt.Fprintf(il.w, "[%d] %s: %s\n", msg.Iter, msg.Level.Fmt(), fmt.Sprintf(msg.Format, msg.Args...))
			if err != nil {
				panic(fmt.Errorf("failed to write log message: [%d] %s: %s\n", msg.Iter, msg.Level.Fmt(), fmt.Sprintf(msg.Format, msg.Args...)))
			}
		}
	}

	return i
}

func (il *IterationLogger) RemoveFork(id uint32) {
	delete(il.fork, id)
}

// AllFmt prints all messages to stdout
func (il *IterationLogger) AllFmt() uint32 {
	return il.fmt(nil)
}

// FilterIterFmt prints all messages for a given iteration id to stdout
func (il *IterationLogger) FilterIterFmt(i int) uint32 {
	return il.fmt(func(msg Msg) bool {
		return msg.Iter == i
	})
}

// logLine is the combined logic for pushing a log message to the log channel.
func (il *IterationLogger) logLine(iter int, level Level, format string, args ...interface{}) {
	if false { // cause vet to treat logLine as a printf wrapper
		_ = fmt.Sprintf(format, args...)
	}

	msg := Msg{
		Iter:   iter,
		Level:  level,
		Format: format,
		Args:   args,
	}

	select {
	case <-il.ctx.Done():
	case il.ch <- msg:
	}
}

// LogPanic pushes a log message to the log channel.
// This method sets the level to [LevelPanic].
// LogPanic is intended to be used after recovering a panic.
func (il *IterationLogger) LogPanic(iter int, recovered any) {
	il.logLine(iter, LevelPanic, "%v", recovered)
}

// LogError pushes a log message to the log channel.
// This method sets the level to [LevelError].
// LogError should be given an error.
func (il *IterationLogger) LogError(iter int, err error) {
	il.logLine(iter, LevelError, "%v", err)
}

// LogWarn pushes a log message to the log channel.
// This method sets the level to [LevelWarn].
// LogWarn will use [fmt] formatting to process the second and up arguments.
func (il *IterationLogger) LogWarn(iter int, args ...interface{}) {
	il.logLine(iter, LevelWarn, "%s", fmt.Sprint(args...))
}

// LogWarnf pushes a log message to the log channel.
// This method sets the level to [LevelWarn].
// LogWarnf will use [fmt] formatting to process the third and up arguments,
// using the format argument.
func (il *IterationLogger) LogWarnf(iter int, format string, args ...interface{}) {
	if false {
		// cause vet to treat LogLinef as a printf wrapper
		_ = fmt.Sprintf(format, args...)
	}

	il.logLine(iter, LevelWarn, "%s", fmt.Sprintf(format, args...))
}

// Log pushes a log message to the log channel.
// This method sets the level to [LevelInfo].
// Log will use [fmt] formatting to process the second and up arguments.
func (il *IterationLogger) Log(iter int, args ...interface{}) {
	il.logLine(iter, LevelInfo, "%s", fmt.Sprint(args...))
}

// Logf pushes a log message to the log channel.
// This method sets the level to [LevelInfo].
// Logf will use [fmt] formatting to process the third and up arguments,
// using the format argument.
func (il *IterationLogger) Logf(iter int, format string, args ...interface{}) {
	if false {
		// cause vet to treat LogLinef as a printf wrapper
		_ = fmt.Sprintf(format, args...)
	}

	il.logLine(iter, LevelInfo, format, args...)
}

// LogDebug pushes a log message to the log channel.
// This method sets the level to [LevelDebug].
// LogDebug will use [fmt] formatting to process the second and up arguments.
func (il *IterationLogger) LogDebug(iter int, args ...interface{}) {
	il.logLine(iter, LevelDebug, "%s", fmt.Sprint(args...))
}

// LogDebugf pushes a log message to the log channel.
// This method sets the level to [LevelDebug].
// LogDebugf will use [fmt] formatting to process the third and up arguments,
// using the format argument.
func (il *IterationLogger) LogDebugf(iter int, format string, args ...interface{}) {
	if false {
		// cause vet to treat LogLinef as a printf wrapper
		_ = fmt.Sprintf(format, args...)
	}

	il.logLine(iter, LevelDebug, format, args...)
}
