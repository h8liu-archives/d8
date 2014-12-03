package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"strings"

	"github.com/h8liu/d8/domain"
)

type Request struct {
	Name     string
	Domains  []string
	Callback string
}

type Server struct{}

func checkIdent(s string) bool {
	if len(s) == 0 {
		return false
	}

	for _, r := range s {
		if r >= 'a' && r <= 'z' {
			continue
		}
		if r >= 'A' && r <= 'Z' {
			continue
		}
		if r >= '0' && r <= '9' {
			continue
		}
		if r == '-' || r == '_' {
			continue
		}
		return false
	}

	return true
}

func checkName(name string) bool {
	p := strings.Index(name, ".")
	if p == -1 {
		return checkIdent(name)
	}

	folder := name[:p]
	file := name[p+1:]

	return checkIdent(folder) && checkIdent(file)
}

func (s *Server) Crawl(req *Request, err *string) error {
	if !checkName(req.Name) {
		*err = "bad job name"
		return nil
	}

	var doms []*domain.Domain

	for _, d := range req.Domains {
		dom, e := domain.Parse(d)
		if e != nil {
			*err = e.Error()
			return nil
		}
		doms = append(doms, dom)
	}

	j := newJob(req.Name, doms, req.Callback)
	go j.run()

	*err = "" // no error
	return nil
}

func serve() {
	s := new(Server)
	e := rpc.Register(s)
	if e != nil {
		log.Fatal(e)
	}
	rpc.HandleHTTP()

	addr := ":5353"
	log.Printf("listening on: %q\n", addr)

	l, e := net.Listen("tcp", addr)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	for {
		e = http.Serve(l, nil)
		if e != nil {
			log.Fatal("serve error:", e)
		}
	}
}

var (
	jobName    = flag.String("o", "", "job output name")
	inputPath  = flag.String("i", "doms", "input domain list")
	serverAddr = flag.String("s", "localhost:5353", "server address")
)

func main() {
	flag.Parse()

	if *jobName == "" {
		serve()
		return
	}

	bs, e := ioutil.ReadFile(*inputPath)
	if e != nil {
		log.Fatal(e)
	}

	req := new(Request)

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
