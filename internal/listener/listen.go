package listener

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"

	"github.com/praveenmahasena/server/internal/pubsub"
)

type Port string
type NetWork string

type Listener struct {
	port    Port
	netWork NetWork
	ctx     context.Context
}

func New(c context.Context, n NetWork, p Port) *Listener {
	return &Listener{
		p,
		n,
		c,
	}
}

func (l *Listener) Run() error {
	lConfig := net.ListenConfig{}
	li, liErr := lConfig.Listen(l.ctx, string(l.netWork), string(l.port))

	if liErr != nil {
		return liErr
	}

	genChat := pubsub.NewGeneralPubSub()
	go genChat.StreamMgs()

	for {
		con, conErr := li.Accept()

		if conErr != nil {
			log.Println(conErr)
			continue
		}

		genChat.Cons.Insert(con)

		go handle(con, genChat)
	}
}

func handle(con net.Conn, c *pubsub.GeneralPubSub) {
	read(con, c)
}

func read(con net.Conn, c *pubsub.GeneralPubSub) {
	s := bufio.NewScanner(con)

	for s.Scan() {
		msg := s.Text()
		c.WriteMsg(msg)
		fmt.Println(msg)
		if msg == "end" {
			c.Cons.Remove(con)
			con.Close()
		}
	}
}
