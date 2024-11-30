package distribute

import (
	"context"
	"golang.org/x/sync/errgroup"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/build"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tools/logger"
	"runtime"
)

func Map[K comparable, V any](
	ctx context.Context,
	il *logger.IterationLogger,
	m map[K]V,
	fn func(ctx context.Context, k K, v V) error,
) error {
	nrWorkers := runtime.GOMAXPROCS(0) * 2

	g, gctx := errgroup.WithContext(ctx)
	g.SetLimit(nrWorkers)

	i := 0

	for k, v := range m {
		k, v := k, v

		g.Go(func() error {
			if build.DEBUG {
				defer func() {
					if r := recover(); r != nil {
						il.LogPanic(i, r)
						panic(r)
					}
				}()

				il.Logf(i, "%v: %v", k, v)
			}

			if err := fn(gctx, k, v); err != nil {
				if build.DEBUG {
					il.LogError(i, err)
				}
				return err
			}

			return nil
		})

		i++
	}

	if err := g.Wait(); err != nil {
		if build.DEBUG {
			il.LogError(-1, err)
		}
		return err
	}

	return nil
}

func Slice[V any](
	ctx context.Context,
	il *logger.IterationLogger,
	m []V,
	fn func(ctx context.Context, i int, v V) error,
) error {
	nrWorkers := runtime.GOMAXPROCS(0) * 2

	g, gctx := errgroup.WithContext(ctx)
	g.SetLimit(nrWorkers)

	for i, v := range m {
		i, v := i, v

		g.Go(func() error {
			if build.DEBUG {
				defer func() {
					if r := recover(); r != nil {
						il.LogPanic(i, r)
						panic(r)
					}
				}()

				il.Logf(i, "%v", v)
			}

			if err := fn(gctx, i, v); err != nil {
				if build.DEBUG {
					il.LogError(i, err)
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
