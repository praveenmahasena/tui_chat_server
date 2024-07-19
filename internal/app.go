package internal

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/praveenmahasena/server/internal/listener"
)

func Start() error {
	ctx, cancel := context.WithCancel(context.Background())

	wg := &sync.WaitGroup{}
	defer wg.Wait()

	defer cancel()

	go func(c context.CancelFunc) {
		defer c()
		cancelCh := make(chan os.Signal, 1)

		signal.Notify(cancelCh, os.Interrupt, syscall.SIGTERM, syscall.SIGINT) // the same thing

		sig := <-cancelCh
		fmt.Printf("Received signal: %s\n", sig)
	}(cancel)

	wg.Add(1)
	go bar(ctx, wg)

	if err := foo(); err != nil {
		return err
	}

	l := listener.New(ctx, "tcp", ":42069")
	return l.Run()
}

func bar(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	<-ctx.Done()
	fmt.Printf("WTF\n")
}

func foo() error {
	return nil
}
