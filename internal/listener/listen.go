package listener

import (
	"bufio"
	"fmt"
	"io"
	"net"
)

type Clients struct {
	con []net.Conn
}

func Run() error {
	li, liErr := net.Listen("tcp", ":42069")
	if liErr != nil {
		return liErr
	}

	for {
		con, conErr := li.Accept()
		clients := Clients{}
		clients.con = append(clients.con, con)

		if conErr != nil {
			continue
		}

		messageCh := make(chan string)
		go read(con, messageCh)
		go write(clients.con, messageCh)
	}
	//return nil
}

func read(con net.Conn, ch chan<- string) {
	s := bufio.NewScanner(con)

	for s.Scan() {
		fmt.Println(s.Text())
		ch <- s.Text()
	}
}

func write(con []net.Conn, ch <-chan string) {

	for message := range ch {

		for _, client := range con {
			io.WriteString(client, message)
		}
	}
}
