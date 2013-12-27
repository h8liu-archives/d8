package main

import (
	"log"
	"os"
)

func noError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func main() {
	c := &Crawler{
		In:    "list",
		Out:   "a.zip",
		Quota: 30,
		Log:   os.Stderr,
	}

	e := c.Crawl()
	noError(e)
}
