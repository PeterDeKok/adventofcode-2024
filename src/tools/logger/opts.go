package logger

import "io"

type Option func(il *IterationLogger)

func WithWriter(w io.Writer) Option {
	return func(il *IterationLogger) {
		il.w = w
	}
}
