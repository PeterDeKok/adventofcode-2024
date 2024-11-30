package ctx

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func CtxWithTrappedCancel() (context.Context, context.CancelFunc) {
	ch := make(chan os.Signal, 10)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		select {
		case <-ctx.Done():
			fmt.Println("Context done")
		case <-ch:
			fmt.Println("Terminating context")
			cancel()
		}
	}()

	signal.Notify(ch, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)

	return ctx, cancel
}
