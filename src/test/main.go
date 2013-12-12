package main

import (
	// "fmt"
	"bufio"
	"log"
	"net"
	"os"
	"strings"

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
	fin, e := os.Open("list")
	noError(e)

	s := bufio.NewScanner(fin)

	for s.Scan() {
		d := D(strings.TrimSpace(s.Text()))
		term.T(tasks.NewInfo(d))
	}
}
