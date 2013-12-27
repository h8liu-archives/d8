package main

import (
	"log"
	"os"
	"runtime"
)

func noError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func main() {
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
