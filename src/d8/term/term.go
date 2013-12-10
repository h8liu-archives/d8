package term

import (
	"io"
	"net"

	"d8/client"
	"d8/domain"
)

type Term struct {
	client *client.Client

	Log       io.Writer
	PrintFlag int
	Retry     int
}

func New(c *client.Client) *Term {
	ret := new(Term)
	ret.client = c
	ret.PrintFlag = client.PrintReply
	ret.Retry = 3
	return ret
}

func (self *Term) Task(t Task) *Branch {
	return newCursor(self).T(t)
}

func (self *Term) Query(d *domain.Domain, t uint16, at net.IP) *Leaf {
	return newCursor(self).Q(d, t, at)
}
