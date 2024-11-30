package distribute

import (
	"context"
	"golang.org/x/sync/errgroup"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/build"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tools/logger"
	"runtime"
	"sync/atomic"
)

func Pool[V any](
	ctx context.Context,
	il *logger.IterationLogger,
	ch <-chan V,
	fn func(ctx context.Context, k int, v V) error,
) (poolSize int, gctx context.Context, g *errgroup.Group) {
	nrWorkers := runtime.GOMAXPROCS(0) * 2

	g, gctx = errgroup.WithContext(ctx)
	g.SetLimit(nrWorkers)

	var i int32
	var n int

	// Start pool nodes until the nrWorkers value is reached
	// or the errgroup defers initialising a new goroutine.
	for ; n < nrWorkers; n++ {
		ok := g.TryGo(func() error {
			select {
			case <-gctx.Done():
				return nil
			case v, ok := <-ch:
				if !ok {
					return nil
				}

				j := atomic.AddInt32(&i, 1)

				if err := fn(gctx, int(j), v); err != nil {
					if build.DEBUG {
						il.LogError(int(j), err)
					}
					return err
				}
			}

			return nil
		})

		if !ok {
			break
		}

		i++
	}

	return n, gctx, g
}
