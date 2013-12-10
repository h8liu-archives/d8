package main

import (
	// "fmt"
	"log"
	"net"

	"d8/client"
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
	term.Q(client.Qs(".", NS, "198.41.0.4"))
	term.Q(client.Qs("liulonnie.net", NS, "74.220.195.131"))
	term.T(tasks.NewRecurType(D("www.yahoo.com"), A))
}
