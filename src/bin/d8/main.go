package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"

	"subcmd"

	"d8/domain"
	"d8/tasks"
	"d8/term"
)

func main() {
	subcmd.Add(console, "", "launch an interactive console")
	subcmd.Add(crawl, "crawl", "crawl a domain list")
	subcmd.Main()
}

func noError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func crawl() {
	runtime.GOMAXPROCS(4)

	c := &Crawler{
		In:    "list",
		Out:   "a.zip",
		Quota: 30,
		Log:   os.Stderr,
	}

	e := c.Crawl()
	noError(e)
}

func console() {
	s := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("d8> ")
		if !s.Scan() {
			break
		}

		line := s.Text()
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		d, e := domain.Parse(line)
		if e != nil {
			fmt.Println("error: ", e)
			continue
		}

		term.T(tasks.NewInfo(d))
		fmt.Println()
	}

	noError(s.Err())

	fmt.Println()
}
