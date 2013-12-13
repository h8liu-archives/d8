package main

import (
	// "fmt"
	"log"
	"net"

	// "d8/client"
	. "d8/domain"
	//. "d8/packet/consts"
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
	term.T(tasks.NewInfo(D("www.peopletopeople.com")))
}
