package jobman

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type job struct {
	name    string
	state   int
	total   int
	crawled int
	sample  string
	err     string
	birth   string
	death   string

	savePath string
	doms     []string
}

func newJob(name string) *job {
	ret := new(job)
	ret.name = name
	return ret
}

func (j *job) Save() error {
	f, e := os.Create(j.savePath)
	if e != nil {
		return e
	}

	for _, d := range j.doms {
		fmt.Fprintln(f, d)
	}

	return f.Close()
}

func (j *job) Load() error {
	f, e := os.Open(j.savePath)
	if e != nil {
		return e
	}

	j.doms = nil
	s := bufio.NewScanner(f)
	for s.Scan() {
		d := s.Text()
		d = strings.TrimSpace(d)
		if d == "" {
			continue
		}

		j.doms = append(j.doms, d)
	}

	e = f.Close()
	if e != nil {
		return e
	}

	return s.Err()
}
