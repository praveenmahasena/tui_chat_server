package main

import (
	"log"

	"github.com/praveenmahasena/server/internal/listener"
)

func main() {

	if err := listener.Run(); err != nil {
		log.Fatalln(err)
	}

}
