package main

import (
	"log"
	"os"
	"runtime"

	"subcmd"
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
	new(Console).Main()
}
