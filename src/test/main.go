package main

import (
	// "fmt"
	"log"
	"net"

	. "d8/domain"
	. "d8/packet/consts"
	"d8/tasks"
	"d8/term"
)

func noError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func main() {
	term.Q(D("."), NS, net.ParseIP("198.41.0.4"))
	term.Q(D("liulonnie.net"), NS, net.ParseIP("74.220.195.131"))
	term.T(tasks.NewRecur(D("liulonnie.net")))
}
