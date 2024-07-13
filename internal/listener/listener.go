package listener

import (
	"context"
	"net"
	"sync"

	"github.com/praveenmahasena/server/internal/pubsub"
)

type Listener struct {
	Port string
}

func New(p string) *Listener {
	return &Listener{
		Port: p,
	}
}

func (l *Listener) Run(ctx context.Context) error {
	wg := &sync.WaitGroup{}

	lc := net.ListenConfig{}

	li, liErr := lc.Listen(ctx, "tcp", l.Port)

	if liErr != nil {
		return liErr
	}

	defer li.Close()

	pub := pubsub.NewPub(wg)

	//go Pub.handleWrite()

	for {
		con, conErr := li.Accept()
		if conErr != nil {
			continue
		}

		pub.Add(con)

		go handle(ctx, con, pub)

	}
	// return nil
}

func handle(ctx context.Context, con net.Conn, p *pubsub.Pub) {
	wg := sync.WaitGroup{}
	wg.Add(1)
	go p.HandleRead(ctx, con)
	wg.Wait()
}
