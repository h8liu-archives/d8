package crawler

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"strings"

	"github.com/h8liu/d8/domain"
)

type Server struct {
	archive string
}

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
	j.archive = s.archive
	go j.run()

	*err = "" // no error
	return nil
}

func Serve(path string) {
	s := rpc.NewServer()
	server := &Server{path}
	e := s.RegisterName("Server", server)
	if e != nil {
		log.Fatal(e)
	}
	rpc.HandleHTTP()

	addr := ":5353"
	log.Printf("listening on: %q\n", addr)

	conn, e := net.Listen("tcp", addr)
	if e != nil {
		log.Fatal("listen error:", e)
	}

	for {
		e = http.Serve(conn, s)
		if e != nil {
			log.Fatal("serve error:", e)
		}
	}
}
