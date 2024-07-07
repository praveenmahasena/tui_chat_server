package pubsub

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
)

type Pub struct {
	pubs   map[net.Conn]bool
	msgQue chan string
	readWg *sync.WaitGroup
	mu     sync.Mutex
}

func NewPub(wg *sync.WaitGroup) *Pub {
	pub := &Pub{
		make(map[net.Conn]bool),
		make(chan string),
		&sync.WaitGroup{},
		sync.Mutex{},
	}

	wg.Add(1)
	go pub.handleWrite(wg)

	return pub
}

func (p *Pub) Add(con net.Conn) {
	p.pubs[con] = true
}

func (p *Pub) HandleRead(ctx context.Context, con net.Conn) {
	p.readWg.Add(1)
	defer p.readWg.Done()

	var data = make([]byte, 1000)

	for {
		var str string
		l, err := con.Read(data)
		if err != nil {
			log.Println(err)
			if ctx.Err() != nil {
				return
			}
		}
		fmt.Printf("received %s from conn %p\n", data, con)
		str = string(data[:l])
		p.msgQue <- str
	}
}

func (p *Pub) handleWrite(wg *sync.WaitGroup) {
	defer wg.Done()

	defer func() {
		for c, _ := range p.pubs {
			c.Write([]byte("server shutting down"))
		}
	}()

	for msg := range p.msgQue {
		for c, _ := range p.pubs {
			fmt.Printf("going to send %s to %p\n", msg, c)
			writeBytes, writeErr := c.Write([]byte(msg))
			if writeErr != nil {
				fmt.Printf("write error to %p: %s\n", c, writeErr)
			} else {
				fmt.Printf("write %d bytes to %p\n", writeBytes, c)
			}
		}
	}
}
