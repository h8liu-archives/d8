package main

import (
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

type Respond struct {
	Name    string
	Crawled int
	Total   int
	Done    bool
	Error   string
}

type server struct{}

func checkName(name string) bool {
	p := strings.Index(name, "/")
	if p == -1 {
		return false
	}

	folder := name[:p]
	file := name[p+1:]
	if len(folder) == 0 {
		return false
	}
	if len(file) == 0 {
		return false
	}

	for _, r := range folder {
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

func (s *server) Crawl(req *Request, err *string) error {
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

func main() {
	s := new(server)
	e := rpc.Register(s)
	if e != nil {
		log.Fatal(e)
	}
	rpc.HandleHTTP()

	l, e := net.Listen("tcp", ":5353")
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
