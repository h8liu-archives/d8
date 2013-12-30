package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"d8/domain"
	"d8/tasks"
	"d8/term"
)

const (
	ModeInfo = iota
	ModeIP
	ModeDig
)

type Console struct {
	Verbose bool
	Mode    int
	Exit    bool
}

func (self *Console) info(doms []string) {
	for _, s := range doms {
		d, e := domain.Parse(s)
		if e != nil {
			fmt.Fprintln(os.Stderr, "error: ", e)
			return
		}

		term.T(tasks.NewInfo(d))
	}
}

func (self *Console) ip(doms []string) {
	for _, s := range doms {
		d, e := domain.Parse(s)
		if e != nil {
			fmt.Fprintln(os.Stderr, "error: ", e)
			return
		}

		term.T(tasks.NewIPs(d))
	}
}

func (self *Console) dig(doms []string) {
	panic("todo")
}

func (self *Console) help() {
	p := func(dot, desc string) { fmt.Printf("%.-10s %s\n", dot, desc) }

	p("help", "print this message")
	p("info", "detect dns structure (default mode)")
	p("ip", "query recursively for ip addresses")
	p("dig", "works like dig")
	p("verbose", "turn on log printing")
	p("quiet", "turn off log printing")
	p("exit", "exit")
}

func (self *Console) line(line string) {
	fields := strings.Fields(line)
	if len(fields) == 0 {
		return
	}

	if fields[0][0] == '.' {
		switch fields[0] {
		case ".verbose":
			fmt.Println("verbose=ture")
			self.Verbose = true
		case ".quiet":
			fmt.Println("verbose=false")
			self.Verbose = false
		case ".info":
			self.Mode = ModeInfo
		case ".dig":
			self.Mode = ModeDig
		case ".ip":
			self.Mode = ModeIP
		case ".help":
			self.help()
		case ".exit":
			self.Exit = true
		default:
			fmt.Fprintln(os.Stderr, "unknown dot command")
		}

		if len(fields) > 1 {
			fmt.Fprintln(os.Stderr, "(other fields ignored)")
		}
		return
	}

	switch self.Mode {
	case ModeInfo:
		self.info(fields)
	case ModeIP:
		self.ip(fields)
	case ModeDig:
		self.dig(fields)
	}
}

func (self *Console) Main() {
	s := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("d8> ")
		if !s.Scan() {
			break
		}

		self.line(strings.TrimSpace(s.Text()))
		if self.Exit {
			break
		}
	}

	noError(s.Err())

	fmt.Println()
}
