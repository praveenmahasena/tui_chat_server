package listener

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"

	"github.com/praveenmahasena/server/internal/pubsub"
)

type Listener struct {
	Port    string
	NetWork string
	Ctx     context.Context
}

func New(c context.Context, n, p string) *Listener {
	return &Listener{
		Port:    p,
		NetWork: n,
		Ctx:     c,
	}
}

func (l *Listener) Run() error {
	lConfig := net.ListenConfig{}
	li, liErr := lConfig.Listen(l.Ctx, l.NetWork, l.Port)

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
