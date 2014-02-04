package main

import (
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/h8liu/d8/subcmd"
	"github.com/h8liu/d8/client"
	"github.com/h8liu/d8/domain"
	"github.com/h8liu/d8/tasks"
	"github.com/h8liu/d8/term"
)

func main() {
	subcmd.Default(single)
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
	new(Console).Main()
}

func single() {
	c, e := client.New()
	noError(e)
	t := term.New(c)
	t.Log = nil
	t.Out = os.Stdout

	for _, s := range os.Args[1:] {
		d, e := domain.Parse(s)
		if e != nil {
			fmt.Fprintln(os.Stderr, e)
			continue
		}
		fmt.Printf("// %v\n", d)

		_, e = t.T(tasks.NewInfo(d))
		if e != nil {
			fmt.Fprintln(os.Stderr, e)
		}
	}
}
