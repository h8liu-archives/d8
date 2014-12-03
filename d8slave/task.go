package main

import (
	"bytes"
	"log"

	"github.com/h8liu/d8/client"
	"github.com/h8liu/d8/domain"
	"github.com/h8liu/d8/tasks"
	"github.com/h8liu/d8/term"
)

type task struct {
	domain *domain.Domain
	client *client.Client
	job    *job
	quota  int

	res string
	out string
	log string
	err string
}

func (task *task) run() {
	defer task.job.taskDone(task)

	logBuf := new(bytes.Buffer)
	t := term.New(task.client)
	t.Log = logBuf

	info := tasks.NewInfo(task.domain)
	_, err := t.T(info)

	if err == nil {
		task.out = info.Out()
		task.res = info.Result()
	} else {
		task.err = err.Error()
	}

	task.log = logBuf.String()

	log.Printf("... %v", task.domain)
}
