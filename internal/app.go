package internal

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/praveenmahasena/server/internal/listener"
)

func Start() error {

	ctx, done := context.WithCancel(context.Background())

	go func() {
		signalChannel := make(chan os.Signal, 1)
		signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)

		select {
		case sig := <-signalChannel:
			fmt.Printf("Received signal: %s\n", sig)
			done()
		case <-ctx.Done():
			fmt.Printf("closing signal goroutine\n")
		}
	}()

	l := listener.New(":42069")
	return l.Run(ctx)

}
