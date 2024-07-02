package listener

import (
	"bufio"
	"io"
	"net"
)

func Run() error {
	li, liErr := net.Listen("tcp", ":42069")
	if liErr != nil {
		return liErr
	}

	clients := []net.Conn{}

	for {
		con, conErr := li.Accept()

		if conErr != nil {
			continue
		}
		clients = append(clients, con)
		messageCh := make(chan string)
		go read(con, messageCh)

		go func() {

			for m := range messageCh {
				for _, c := range clients {
					io.WriteString(c, m)
				}
			}

		}()

	}
	//return nil
}

func read(con net.Conn, ch chan<- string) {
	s := bufio.NewScanner(con)

	for s.Scan() {
		msg := s.Text()
		if msg == "end" {
			ch <- con.RemoteAddr().String() + "left"
			con.Close()
			continue
		}
		ch <- msg
	}
}
