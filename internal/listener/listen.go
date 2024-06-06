package listener

import (
	"fmt"
	"net"
)

type Listener struct {
	Network, Addr string
}

func New(n, a string) *Listener {
	return &Listener{
		Network: n,
		Addr:    a,
	}
}

func (l *Listener) Listen() error {
	li, liErr := net.Listen(l.Network, l.Addr)

	if liErr != nil {
		return liErr
	}

	for {
		con, conErr := li.Accept()
		if conErr != nil {
			continue
		}
		fmt.Println(con)
		con.Write([]byte("connected to server"))
	}

	//return nil
}
