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

func ip(s string) net.IP {
	return net.ParseIP(s)
}

func main() {
	term.Q(D("."), NS, ip("198.41.0.4"))
	term.Q(D("liulonnie.net"), NS, ip("74.220.195.131"))
	term.T(tasks.NewRecurType(D("www.yahoo.com"), A))
}
