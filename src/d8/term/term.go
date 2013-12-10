package term

import (
	"fmt"
	"io"
	"os"

	"d8/client"
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

	ret, e := newCursor(self).T(t)
	if e != nil {
		panic(e)
	}

	self.done++

	return ret
}

func (self *Term) Q(q *client.Query) *Leaf {
	if self.done != 0 {
		fmt.Fprintln(self.Log)
	}

	ret, e := newCursor(self).Q(q)
	if e != nil {
		panic(e)
	}

	self.done++

	return ret
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

func Q(q *client.Query) *Leaf {
	return makeStd().Q(q)
}
