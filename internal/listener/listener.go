package listener

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
)

type Listner struct {
	Port string
}

func New(port string) *Listner {
	return &Listner{
		Port: port,
	}
}

type Pub struct {
	pub    map[string][]net.Conn
	msgQue chan Msg
}

func (l *Listner) Run() error {
	li, liErr := net.Listen("tcp", l.Port)

	if liErr != nil {
		return liErr
	}

	p := Pub{
		pub:    make(map[string][]net.Conn),
		msgQue: make(chan Msg),
	}
	go p.Watch()
	go p.writeToAll()

	for {
		con, conErr := li.Accept()
		if conErr != nil {
			log.Println(conErr)
			continue
		}
		msg := Msg{}
		json.NewDecoder(con).Decode(&msg)
		if p.pub[msg.ChatType] == nil {
			p.pub[msg.ChatType] = []net.Conn{con}
		} else {
			p.pub[msg.ChatType] = append(p.pub[msg.ChatType], con)
		}
		go handleCon(p, con)
	}

	// return nil
}

// TODO : missing 1st message

type Msg struct {
	UserName string
	ChatType string
	// Password string
	Message string
}

func handleCon(p Pub, con net.Conn) {
	r := bufio.NewReader(con)
	msg := Msg{}

	for {
		msgB, err := r.ReadBytes('\n')
		if err != nil {
			continue
		}
		json.Unmarshal(msgB, &msg)
		p.msgQue <- msg
	}
}

func (p *Pub) Watch() {
	for msg := range p.msgQue {
		go send(p, msg)
	}

}

func send(p *Pub, msg Msg) {
	for _, c := range p.pub[msg.ChatType] {
		go func() {
			m, _ := json.Marshal(msg)
			n, e := c.Write(m)
			fmt.Println(n, e)
		}()
	}
}

func (p *Pub) writeToAll() {
	s := bufio.NewScanner(os.Stdin)

	for s.Scan() {
		for _, cons := range p.pub {
			for _, c := range cons {
				go func(con net.Conn) {
					fmt.Println(con.Write(s.Bytes()))
				}(c)
			}
		}
	}
}
