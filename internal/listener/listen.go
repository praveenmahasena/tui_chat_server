package listener

import (
	"bufio"
	"log"
	"net"
	"os"
	"sync"
)

type Listner struct {
	Network string
	Addr    string
}

func New(n, a string) *Listner {
	return &Listner{
		Network: n,
		Addr:    a,
	}
}

func (l *Listner) Listen() error {
	li, liErr := net.Listen(l.Network, l.Addr)

	if liErr != nil {
		return liErr
	}

	defer li.Close()

	go handleConnect(li)
	ch := make(chan struct{}) // at this time I would be using sync.WaitGroup{} but for practice I'm using empty struct{} channel

	<-ch

	return nil
}

func handleConnect(li net.Listener) {
	for {
		con, conErr := li.Accept()
		if conErr != nil {
			continue
		}

		go readWrite(con)
	}

}

func readWrite(con net.Conn) {
	wg := sync.WaitGroup{}
	wg.Add(2)
	go read(con, &wg)
	go write(con, &wg)
	wg.Wait()
	defer con.Close()
}

func write(con net.Conn, wg *sync.WaitGroup) {
	s := bufio.NewScanner(os.Stdin)

	for s.Scan() {
		con.Write(s.Bytes())
		if s.Text() == "end" {
			wg.Done()
			break
		}
	}
}

func read(con net.Conn, wg *sync.WaitGroup) {
	s := bufio.NewScanner(con)
	for s.Scan() {
		log.Println(s.Text())
		if s.Text() == "end" {
			wg.Done()
			break
		}
	}
}
