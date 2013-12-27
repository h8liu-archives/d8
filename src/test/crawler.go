package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"sync"

	"d8/client"
	. "d8/domain"
	"d8/packet/rdata"
	"d8/tasks"
	"d8/term"
	"printer"
)

type Crawler struct {
	In    string
	Out   string
	Quota int
	Log   io.Writer
}

type crawlTask struct {
	id     int
	domain *Domain
	client *client.Client
	out    *zip.Writer
	lock   *sync.Mutex

	quota       int
	quotaReturn chan int
}

func (self *Crawler) quotas() chan int {
	nquota := self.Quota
	if nquota == 0 {
		nquota = 1
	}
	ret := make(chan int, nquota)
	for i := 0; i < nquota; i++ {
		ret <- i
	}

	return ret
}

func (self *Crawler) Crawl() error {
	// load input
	list, e := LoadList(self.In, self.Log)
	if e != nil {
		return e
	}

	// init output
	fout, e := os.Create(self.Out)
	if e != nil {
		return e
	}
	out := zip.NewWriter(fout)
	defer out.Close()

	c, e := client.New()
	if e != nil {
		return e
	}

	quotas := self.quotas()
	lock := new(sync.Mutex)

	for id, d := range list {
		q := <-quotas
		task := &crawlTask{
			id:          id,
			domain:      d,
			client:      c,
			out:         out,
			lock:        lock,
			quota:       q,
			quotaReturn: quotas,
		}

		go task.run()
	}

	// join
	for i := 0; i < cap(quotas); i++ {
		<-quotas
	}

	return nil
}

func (self *crawlTask) create(path string) (io.Writer, error) {
	return self.out.CreateHeader(&zip.FileHeader{
		Name:   path,
		Method: zip.Store,
	})
}

func (self *crawlTask) run() {
	logbuf := new(bytes.Buffer)
	t := term.New(self.client)
	t.Log = logbuf
	info := tasks.NewInfo(self.domain)
	_, err := t.T(info)

	d := self.domain
	self.lock.Lock()

	fout, e := self.create(fmt.Sprintf("log/%s", d))
	noError(e)
	_, e = io.Copy(fout, logbuf)
	noError(e)

	fout, e = self.create(fmt.Sprintf("out/%s", d))
	if err == nil {
		e = printInfo(info, fout)
		noError(e)
	} else {
		fmt.Fprintf(fout, "error: %v\n", err)
	}

	self.lock.Unlock()

	self.quotaReturn <- self.quota
}

func printInfo(info *tasks.Info, out io.Writer) error {
	p := printer.New(out)

	p.Printf("%v {", info.Domain)
	p.ShiftIn()

	if len(info.Cnames) > 0 {
		p.Print("cnames {")
		p.ShiftIn()
		for _, r := range info.Cnames {
			p.Printf("%v -> %v", r.Domain, rdata.ToDomain(r.Rdata))
		}
		p.ShiftOut()
		p.Print("}")
	}

	if len(info.Results) == 0 {
		p.Print("(unresolvable)")
	} else {
		p.Print("ips {")
		p.ShiftIn()

		for _, r := range info.Results {
			d := r.Domain
			ip := rdata.ToIPv4(r.Rdata)
			if d.Equal(info.Domain) {
				p.Printf("%v", ip)
			} else {
				p.Printf("%v(%v)", ip, d)
			}
		}

		p.ShiftOut()
		p.Print("}")
	}

	if len(info.NameServers) > 0 {
		p.Print("servers {")
		p.ShiftIn()

		for _, ns := range info.NameServers {
			p.Printf("%v", ns)
		}

		p.ShiftOut()
		p.Print("}")
	}

	if len(info.Records) > 0 {
		p.Print("records {")
		p.ShiftIn()

		for _, rr := range info.Records {
			p.Printf("%v", rr.Digest())
		}

		p.ShiftOut()
		p.Print("}")
	}

	p.ShiftOut()
	p.Print("}")

	return p.Error
}
