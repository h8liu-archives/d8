package term

import (
	"bytes"
	"io"
	"net"

	"d8/client"
	"d8/domain"
	"printer"
)

type Cursor interface {
	printer.Interface

	T(t Task) *Branch
	Q(d *domain.Domain, t uint16, at net.IP) *Leaf
}

type cursor struct {
	*printer.Printer
	*Term // conveniently inherits the term options
	*stack

	log *bytes.Buffer
}

var _ Cursor = new(cursor)

func newCursor(t *Term) *cursor {
	ret := new(cursor)

	ret.Term = t
	ret.log = new(bytes.Buffer)
	ret.stack = newStack()

	if t.Log == nil {
		ret.Printer = printer.New(ret.log)
	} else {
		ret.Printer = printer.New(io.MultiWriter(t.Log, ret.log))
	}

	return ret
}

func (self *cursor) Q(d *domain.Domain, t uint16, at net.IP) *Leaf {
	ret := self.q(d, t, at)
	self.TopAdd(ret)
	return ret
}

func (self *cursor) T(t Task) *Branch {
	ret := newBranch(t)
	self.TopAdd(ret)
	self.Push(ret)

	t.Run(self)

	b := self.Pop()
	if b != ret || b.Task != t {
		panic("bug")
	}

	return ret
}

func (self *cursor) q(d *domain.Domain, t uint16, at net.IP) *Leaf {
	q := &client.Query{
		Domain: d,
		Type:   t,
		Server: &net.UDPAddr{
			IP:   at,
			Port: client.DNSPort,
		},
		Printer:   self.Printer,
		PrintFlag: self.PrintFlag,
	}

	ret := new(Leaf)
	ret.Attempts = make([]*client.Exchange, 0, self.Retry)

	for i := 0; i < self.Retry; i++ {
		answer := self.client.Query(q)
		ret.Attempts = append(ret.Attempts, answer)
		if answer.Timeout() {
			self.Print("// retry")
			continue
		}
		break
	}

	return ret
}
