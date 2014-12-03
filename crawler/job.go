package crawler

import (
	"database/sql"
	"log"
	"net/rpc"

	_ "github.com/mattn/go-sqlite3"

	"github.com/h8liu/d8/client"
	"github.com/h8liu/d8/domain"
)

type Respond struct {
	Name    string
	Crawled int
	Total   int
	Done    bool
	Error   string
}

type job struct {
	name     string
	domains  []*domain.Domain
	callback string
	crawled  int
	respond  *Respond

	db     *sql.DB
	client *rpc.Client
	quotas chan int

	resChan   chan *task
	writeDone chan bool
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

	if e != nil {
		log.Print(j.name, e)
	}
}

func (j *job) cleanup() {
	if j.db != nil {
		e := j.db.Close()
		if e != nil {
			log.Print(j.db, e)
		}
	}

	if j.client != nil {
		e := j.client.Close()
		if e != nil {
			log.Print(j.name, e)
		}
	}
}

func (j *job) fail(e error) {
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

func (j *job) taskDone(t *task) {
	j.resChan <- t
	j.quotas <- t.quota
}

func (j *job) run() {
	log.Printf("job %s started", j.name)
	defer log.Printf("job %s done", j.name)
	defer j.cleanup()

	db, err := sql.Open("sqlite3", j.name)
	if j.failOn(err) {
		return
	}

	j.db = db

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

	if !q(`create table jobs (
			domain text not null primary key,
			output text not null,
			result text not null,
			err text not null,
			log text not null)`) {
		return
	}

	log.Printf("job %s starts crawling", j.name)
	j.crawl()
}

func (j *job) makeQuotas() chan int {
	nquota := 300
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

	j.quotas = j.makeQuotas()
	j.resChan = make(chan *task, 300)
	defer close(j.resChan)

	// launch the jobs
	go func() {
		for _, d := range j.domains {
			quota := <-j.quotas
			task := &task{
				domain: d,
				client: c,
				job:    j,
				quota:  quota,
			}
			go task.run()
		}
	}()

	j.writeOut()
}

func (j *job) writeOut() {
	n := 0
	total := j.respond.Total

	chkerr := func(e error) bool {
		if e != nil {
			j.respond.Error = e.Error()
			j.cb()
			return true
		}
		return false
	}

	const insertStmt = `insert into jobs
		(domain, output, result, err, log) values
		(?, ?, ?, ?, ?)`

	tx, err := j.db.Begin()
	if chkerr(err) {
		return
	}
	stmt, err := tx.Prepare(insertStmt)
	if chkerr(err) {
		return
	}

	for n < total {
		t := <-j.resChan

		_, err = stmt.Exec(t.domain.String(),
			t.out, t.res, t.err, t.log,
		)
		if chkerr(err) {
			return
		}

		n++
		if n%5000 == 0 {
			err = tx.Commit()
			if chkerr(err) {
				return
			}

			tx, err = j.db.Begin()
			if chkerr(err) {
				return
			}
			stmt, err = tx.Prepare(insertStmt)
			if chkerr(err) {
				return
			}

			j.respond.Crawled = n
			j.cb()
		}
	}

	err = tx.Commit()
	if chkerr(err) {
		return
	}

	j.respond.Done = true
	j.cb()
}
