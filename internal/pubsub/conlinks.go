package pubsub

import (
	"log"
	"net"
)

type ConNode struct {
	prev *ConNode
	con  net.Conn
	next *ConNode
}

type ConList struct {
	head *ConNode
	tail *ConNode
}

func NewConList() *ConList {
	return &ConList{}
}

func (c *ConList) Insert(con net.Conn) bool {
	conNode := &ConNode{
		con: con,
	}

	if c.head == nil {
		c.head, c.tail = conNode, conNode
		return true
	}
	c.tail.next, c.tail, conNode.prev = conNode, conNode, c.tail
	return true
}

func (c *ConList) Remove(con net.Conn) bool {
	// here
	if c.head.con == con {
		if c.head.next == nil {
			c.head = nil
		} else {
			c.head, c.head.next.prev, c.head.next = c.head.next, nil, nil
		}
		return true
	}

	if c.tail.con == con {
		c.tail, c.tail.prev, c.tail.prev.next = c.tail.prev, nil, nil
		return true
	}

	n := c.head

	for n.next != nil {
		if n.con == con {
			n.prev.next = n.next
			return true
		}
		n = n.next
	}

	return true
}

func (c *ConList) Write(msg string) {
	if c.head == nil {
		return
	}
	n := c.head
	for n.next != nil {
		_, err := n.con.Write([]byte(msg))
		if err != nil {
			log.Println(err)
			continue
		}
		n = n.next
	}
}
