package distribute

import (
	"context"
	"golang.org/x/sync/errgroup"
	"io"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/build"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tools/input"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tools/logger"
	"runtime"
)

func Lines(ctx context.Context, il *logger.IterationLogger, rd io.Reader, fn func(ctx context.Context, i int, l string) error) error {
	nrWorkers := runtime.GOMAXPROCS(0) * 2

	g, gctx := errgroup.WithContext(ctx)
	g.SetLimit(nrWorkers)

	for lineNr, line := range input.LineReader(rd) {
		lineNr, line := lineNr, line

		g.Go(func() error {
			if build.DEBUG {
				defer func() {
					if r := recover(); r != nil {
						il.LogPanic(lineNr, r)
						panic(r)
					}
				}()

				il.Logf(lineNr, "%s", line)
			}

			if err := fn(gctx, lineNr, line); err != nil {
				if build.DEBUG {
					il.LogError(lineNr, err)
				}
				return err
			}

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		if build.DEBUG {
			il.LogError(-1, err)
		}
		return err
	}

	return nil
}
