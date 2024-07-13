package listener

import (
	"log"
	"net"
	"sync"
)

type Listener struct {
	Port string
}

func New(p string) *Listener {
	return &Listener{
		Port: p,
	}
}

type Pub struct {
	pubs   map[net.Conn]bool
	msgQue chan string
	mu     sync.Mutex
}

func (p *Pub) Add(con net.Conn) {
	p.pubs[con] = true
}

func (l *Listener) Run() error {
	li, liErr := net.Listen("tcp", l.Port)

	if liErr != nil {
		return liErr
	}

	defer li.Close()

	Pub := &Pub{
		pubs:   map[net.Conn]bool{},
		msgQue: make(chan string),
		mu:     sync.Mutex{},
	}

	//go Pub.handleWrite()

	for {
		con, conErr := li.Accept()
		if conErr != nil {
			continue
		}

		Pub.Add(con)

		go handle(con, Pub)

	}
	// return nil
}

func handle(con net.Conn, p *Pub) {
	wg := sync.WaitGroup{}
	wg.Add(1)
	go handleRead(con, p)
	wg.Wait()
}

func (p *Pub) AddQue(str string) {
	p.mu.Lock()
	p.msgQue <- str
	p.mu.Unlock()
}

func handleRead(con net.Conn, p *Pub) {
	var data = make([]byte, 1000)

	for {
		var str string
		l, err := con.Read(data)
		if err != nil {
			log.Println(err)
		}
		str = string(data[:l])

		p.AddQue(str)

	}
}

func (p *Pub) handleWrite() {
	for msg := range p.msgQue {
		for c, _ := range p.pubs {
			c.Write([]byte(msg))
		}
	}
}
