package distribute

import (
	"context"
	"golang.org/x/sync/errgroup"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/build"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tools/logger"
	"runtime"
)

func Group(ctx context.Context, il *logger.IterationLogger, fn func(g *errgroup.Group, gctx context.Context) error) error {
	nrWorkers := runtime.GOMAXPROCS(0) * 2

	cctx, cancel := context.WithCancelCause(ctx)
	g, gctx := errgroup.WithContext(cctx)
	g.SetLimit(nrWorkers)

	if err := fn(g, gctx); err != nil {
		if build.DEBUG {
			il.LogError(-1, err)
		}
		cancel(err)
	}

	if err := g.Wait(); err != nil {
		if build.DEBUG {
			il.LogError(-1, err)
		}
		return err
	}

	return nil
}
