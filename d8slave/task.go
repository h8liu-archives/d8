package main

import (
	"bytes"
	"sync"

	"github.com/h8liu/d8/client"
	"github.com/h8liu/d8/domain"
	"github.com/h8liu/d8/printer"
	"github.com/h8liu/d8/tasks"
	"github.com/h8liu/d8/term"
)

type task struct {
	domain *domain.Domain
	client *client.Client
	wait   *sync.WaitGroup
}

func (task *task) run() {
	defer task.wait.Done()

	logBuf := new(bytes.Buffer)
	t := term.New(task.client)
	t.Log = logBuf

	outBuf := new(bytes.Buffer)

	info := tasks.NewInfo(task.domain)
	_, err := t.T(info)

	if err == nil {
		p := printer.New(outBuf)
		info.PrintTo(p)
	}

	panic("save log, out, err, and result")

}
