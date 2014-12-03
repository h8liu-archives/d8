// Package jobman defines the job managing server
package jobman

import (
	"bytes"
	"database/sql"
	"errors"
	"log"
	// "net/rpc"
	"fmt"
	"strings"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// JobMan defines a crawling job manager
type JobMan struct {
	dbpath string
	db     *sql.DB

	jobs    map[string]*job
	lock    *sync.Mutex
	newJobs chan *job

	workers []*worker // worker client addresses
}

// NewJobMan creates a new job manager
func NewJobMan(dbpath string) (*JobMan, error) {
	ret := new(JobMan)
	ret.dbpath = dbpath

	var err error
	ret.db, err = sql.Open("sqlit3", dbpath)
	if err != nil {
		return nil, err
	}

	ret.lock = new(sync.Mutex)
	ret.newJobs = make(chan *job, 10)

	return ret, nil
}

// AddWorker adds a worker address to the job man for handling jobs
func (jm *JobMan) AddWorker(addr string) {
	w := newWorker(addr)
	jm.workers = append(jm.workers, w)
}

// Serve starts crawling jobs.
// Blcoking call, will never return.
func (jm *JobMan) Serve() error {
	if len(jm.workers) == 0 {
		return errors.New("no workers")
	}

	go jm.serveCallback()

	for {
		// here is the scheduler
		// listen to new jobs and progress call backs
		// web queries will just goto the database directly
		panic("todo")
	}
}

func (jm *JobMan) q(sql string, args ...interface{}) {
	_, e := jm.db.Exec(sql, args...)
	if e != nil {
		log.Printf("sql: %s\n", sql)
		log.Fatal(e)
	}
}

// makeSample creates a sample for the domain list.
// currently it just looks at the head.
// it serves as a description.
func makeSample(doms []string) string {
	ret := new(bytes.Buffer)
	cnt := 0
	for _, d := range doms {
		d = strings.TrimSpace(d)
		if d == "" {
			continue
		}
		fmt.Fprintf(ret, "%s\n", d)
		cnt++
		if cnt >= 20 {
			break
		}
	}

	if cnt == 20 {
		fmt.Fprintf(ret, "...\n")
	}

	fmt.Fprintf(ret, "(%d domains)\n", len(doms))

	return ret.String()
}

const (
	JobCreated int = iota
	JobReady
	JobWorking
	JobDone
)

func (jm *JobMan) createJob(name string, doms []string) (*job, error) {
	jm.lock.Lock()
	defer jm.lock.Unlock()

	if jm.jobs[name] != nil {
		return nil, fmt.Errorf("job %q exists", name)
	}

	j := newJob(name)
	j.savePath = name
	e := j.Save()

	if e != nil {
		return nil, e
	}

	j.sample = makeSample(doms)
	j.birth = time.Now().String()
	j.state = JobCreated
	j.doms = doms
	j.total = len(doms)

	jm.jobs[name] = j
	return j, nil
}

// CreateJob creates a new job.
func (jm *JobMan) CreateJob(name string, doms []string) error {
	j, e := jm.createJob(name, doms)
	if e != nil {
		return e
	}

	// make this persistent
	jm.q(`insert into jobs 
		(name, state, total, crawled, sample, birth) 
		values (?, ?, ?, ?, ?, ?)`,
		name, j.state, j.total, 0, j.sample, j.birth,
	)

	jm.newJobs <- j

	return nil
}

func (jm *JobMan) serveCallback() error {
	panic("todo")
}
