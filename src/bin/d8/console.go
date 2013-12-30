package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"d8/client"
	"d8/domain"
	"d8/tasks"
	"d8/term"
)

const (
	ModeInfo = iota
	ModeIP
	ModeRecur
	ModeDig
)

type Console struct {
	Verbose bool
	Mode    int
	Exit    bool
	Term    *term.Term
}

func (self *Console) printError(e error) {
	if e == nil {
		return
	}
	fmt.Fprintln(os.Stderr, "error: ", e)
}

func (self *Console) info(doms []string) {
	for _, s := range doms {
		d, e := domain.Parse(s)
		if e != nil {
			self.printError(e)
			continue
		}

		_, e = self.Term.T(tasks.NewInfo(d))
		self.printError(e)
	}
}

func (self *Console) ip(doms []string) {
	for _, s := range doms {
		d, e := domain.Parse(s)
		if e != nil {
			self.printError(e)
			continue
		}

		_, e = self.Term.T(tasks.NewIPs(d))
		self.printError(e)
	}
}

func (self *Console) dig(doms []string) {
	fmt.Println("(mode not implemented)")
}

func (self *Console) recur(doms []string) {
	fmt.Println("(mode not implemented)")
}

func (self *Console) help() {
	p := func(dot, desc string) { fmt.Printf(".%-10s %s\n", dot, desc) }

	p("info", "ips and detect dns structure (default mode)")
	p("ip", "query recursively for ip addresses")
	p("recur", "works like a recursive dig")
	p("dig", "works like dig")
	p("verbose", "turn on log printing")
	p("quiet", "turn off log printing")
	p("help", "print this message")
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
			self.Term.Log = os.Stdout
		case ".quiet":
			self.Term.Log = nil
		case ".info":
			self.Mode = ModeInfo
		case ".ip":
			self.Mode = ModeIP
		case ".recur":
			self.Mode = ModeRecur
		case ".dig":
			self.Mode = ModeDig
		case ".help":
			self.help()
		case ".exit":
			self.Exit = true
		default:
			self.printError(errors.New("invalid dot command"))
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
	case ModeRecur:
		self.recur(fields)
	case ModeDig:
		self.dig(fields)
	default:
		panic("bug")
	}
}

func (self *Console) Main() {
	s := bufio.NewScanner(os.Stdin)
	if self.Term == nil {
		c, e := client.New()
		noError(e)
		self.Term = term.New(c)
		self.Term.Log = nil
		self.Term.Out = os.Stdout
	}

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
