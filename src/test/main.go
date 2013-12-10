package main

import (
	// "fmt"
	"log"
	"net"
	"os"

	"d8/client"
	. "d8/domain"
	. "d8/packet/consts"
	"d8/term"
)

func noError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func main() {
	c, e := client.New()
	noError(e)

	t := term.New(c)
	t.Log = os.Stdout

	t.Query(D("."), NS, net.ParseIP("198.41.0.4"))
	t.Query(D("liulonnie.net"), NS, net.ParseIP("74.220.195.131"))
}
