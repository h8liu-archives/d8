package term

import (
	"fmt"
	"io"
	"net"
	"os"

	"d8/client"
	"d8/domain"
)

type Term struct {
	client *client.Client

	done      int
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

func (self *Term) T(t Task) *Branch {
	if self.done != 0 {
		fmt.Fprintln(self.Log)
	}

	b, e := newCursor(self).T(t)
	if e != nil {
		panic(e)
	}

	self.done++

	return b
}

func (self *Term) Q(d *domain.Domain, t uint16, at net.IP) *Leaf {
	if self.done != 0 {
		fmt.Fprintln(self.Log)
	}

	q, e := newCursor(self).Q(d, t, at)
	if e != nil {
		panic(e)
	}

	self.done++

	return q
}

var std *Term

func makeStd() *Term {
	if std == nil {
		c, e := client.New()
		if e != nil {
			panic(e)
		}

		std = New(c)
		std.Log = os.Stdout
	}

	return std
}

func T(t Task) *Branch {
	return makeStd().T(t)
}

func Q(d *domain.Domain, t uint16, at net.IP) *Leaf {
	return makeStd().Q(d, t, at)
}
