package main

import (
	"database/sql"
	"io/ioutil"
	"log"
	"net/rpc"
	"sync"

	_ "github.com/mattn/go-sqlite3"

	"github.com/h8liu/d8/client"
	"github.com/h8liu/d8/domain"
)

type job struct {
	name     string
	domains  []*domain.Domain
	callback string
	crawled  int
	respond  *Respond
	err      error
	db       *sql.DB
	client   *rpc.Client
}

func newJob(name string, doms []*domain.Domain, cb string) *job {
	ret := new(job)
	ret.name = name
	ret.domains = doms
	ret.callback = cb

	ret.respond = new(Respond)
	ret.respond.Name = name
	ret.respond.Total = len(doms)

	return ret
}

func (j *job) connect() error {
	var e error
	j.client, e = rpc.DialHTTP("tcp", j.callback)
	return e
}

func (j *job) call() error {
	var s string
	return j.client.Call("Respond", j.respond, &s)
}

func (j *job) cb() {
	if j.callback == "" {
		return
	}

	// in callback, we log error
	var e error
	if j.client == nil {
		e = j.connect()
		if e != nil {
			log.Print(j.name, e)
			return
		}
	}

	e = j.call()
	if e == rpc.ErrShutdown {
		j.client.Close()
		e = j.connect()
		if e != nil {
			log.Print(j.name, e)
			return
		}
		e = j.call()
	}
}

func (j *job) closeClient() {
	if j.client != nil {
		e := j.client.Close()
		if e != nil && e != rpc.ErrShutdown {
			log.Print(j.name, e)
		}
	}
}

func (j *job) fail(e error) {
	j.err = e
	j.respond.Error = e.Error()
	j.cb()
}

func (j *job) failOn(e error) bool {
	if e != nil {
		j.fail(e)
		return true
	}
	return false
}

func (j *job) quotas() chan int {
	nquota := 300
	if nquota == 0 {
		nquota = 1
	}

	ret := make(chan int, nquota)
	for i := 0; i < nquota; i++ {
		ret <- i
	}

	return ret
}

func (j *job) crawl() {
	c, e := client.New()
	if j.failOn(e) {
		return
	}

	wg := new(sync.WaitGroup)
	wg.Add(300)

	for _, d := range j.domains {
		wg.Wait()
		task := &task{
			domain: d,
			client: c,
			wait:   wg,
		}
		go task.run()
	}
}

func (j *job) run() {
	defer j.closeClient()

	tmp, e := ioutil.TempFile("d8", "job")
	if j.failOn(e) {
		return
	}

	f := tmp.Name()
	if j.failOn(tmp.Close()) {
		return
	}

	db, err := sql.Open("sqlite3", f)
	if j.failOn(err) {
		return
	}

	q := func(sql string) bool {
		_, e := db.Exec(sql)
		if e != nil {
			log.Printf("sql fail: %s\n", sql)
			log.Print(e)
			j.fail(e)
			return false
		}

		return true
	}

	if !q(`create table jobs {
			domain text not null primary key,
			output text not null,
			result text not null,
			err text not null,
			log text not null)`) {
		return
	}

	j.db = db

	j.crawl()
}
