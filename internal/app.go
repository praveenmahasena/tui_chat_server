package internal

import (
	"github.com/praveenmahasena/server/internal/listener"
)

func Start() error {
	l := listener.New("tcp", ":42069")
	return l.Listen()
}
