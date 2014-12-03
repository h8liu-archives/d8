package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/rpc"
	"strings"

	"github.com/h8liu/d8/crawler"
)

var (
	jobName    = flag.String("o", "", "job output name")
	inputPath  = flag.String("i", "doms", "input domain list")
	serverAddr = flag.String("s", "localhost:5353", "server address")
	saveAddr   = flag.String("a", "", "archive prefix")
)

func main() {
	flag.Parse()

	if *jobName == "" {
		crawler.Serve(*saveAddr)
		return
	}

	bs, e := ioutil.ReadFile(*inputPath)
	if e != nil {
		log.Fatal(e)
	}

	req := new(crawler.Request)

	doms := strings.Split(string(bs), "\n")
	for _, d := range doms {
		d = strings.TrimSpace(d)
		if d == "" {
			continue
		}
		req.Domains = append(req.Domains, d)
	}

	req.Name = *jobName

	c, e := rpc.DialHTTP("tcp", *serverAddr)
	if e != nil {
		log.Fatal(e)
	}

	var reply string
	e = c.Call("Server.Crawl", req, &reply)
	if e != nil {
		log.Fatal(e)
	} else if reply != "" {
		log.Print(reply)
	}

	e = c.Close()
	if e != nil {
		log.Fatal(e)
	}
}
