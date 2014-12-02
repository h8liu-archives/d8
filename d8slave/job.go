package main

import (
	"database/sql"
	"io/ioutil"
	"log"

	"github.com/h8liu/d8/domain"
	_ "github.com/mattn/go-sqlite3"
)

type job struct {
	name     string
	domains  []*domain.Domain
	callback string
	crawled  int
	respond  *Respond
	err      error
	db       *sql.DB
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

func (j *job) cb() {
	panic("todo")
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

func (j *job) run() {
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
			log text not null)`) {
		return
	}

	j.db = db
}
