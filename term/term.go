package term

import (
	// "fmt"
	"io"
	"os"

	"github.com/h8liu/d8/client"
	"github.com/h8liu/d8/printer"
)

type Term struct {
	client *client.Client

	done      int
	Log       io.Writer
	Out       io.Writer
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

func (self *Term) T(t Task) (*Branch, error) {
	ret, e := newCursor(self).T(t)
	self.done++

	if e == nil {
		p := printer.New(self.Out)
		t.PrintTo(p)
	}

	return ret, e
}

func (self *Term) Q(q *client.Query) (*Leaf, error) {
	ret, e := newCursor(self).Q(q)
	self.done++

	return ret, e
}

func (self *Term) Count() int {
	return self.done
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
		std.Out = os.Stdout
	}

	return std
}

func T(t Task) *Branch {
	ret, e := makeStd().T(t)
	if e != nil {
		panic(e)
	}
	return ret
}

func Q(q *client.Query) *Leaf {
	ret, e := makeStd().Q(q)
	if e != nil {
		panic(e)
	}
	return ret
}
