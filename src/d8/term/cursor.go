package term

import (
	"errors"
	"net"

	"d8/client"
	"d8/domain"
	"printer"
)

type Cursor interface {
	printer.Interface

	T(t Task) (*Branch, error)
	Q(d *domain.Domain, t uint16, at net.IP) (*Leaf, error)
}

type cursor struct {
	*printer.Printer
	*Term // conveniently inherits the term options
	*stack
	nquery int
}

var _ Cursor = new(cursor)

func newCursor(t *Term) *cursor {
	ret := new(cursor)

	ret.Term = t
	ret.stack = newStack()
	ret.Printer = printer.New(t.Log)

	return ret
}

const (
	MaxDepth = 10
	MaxQuery = 100
)

var (
	errTooDeep        = errors.New("too deep")
	errTooManyQueries = errors.New("too many queries")
)

func (self *cursor) Q(d *domain.Domain, t uint16, at net.IP) (*Leaf, error) {
	if self.nquery >= MaxQuery {
		return nil, errTooManyQueries
	}

	self.nquery++
	ret := self.q(d, t, at)
	self.TopAdd(ret)
	return ret, nil
}

func (self *cursor) T(t Task) (*Branch, error) {
	if self.Len() >= MaxDepth {
		return nil, errTooDeep
	}

	ret := newBranch(t)
	self.TopAdd(ret)
	self.Push(ret)

	t.Run(self)

	b := self.Pop()
	if b != ret || b.Task != t {
		panic("bug")
	}

	return ret, nil
}

func (self *cursor) q(d *domain.Domain, t uint16, at net.IP) *Leaf {
	q := &client.QueryPrinter{
		Query: &client.Query{
			Domain: d,
			Type:   t,
			Server: &net.UDPAddr{
				IP:   at,
				Port: client.DNSPort,
			},
		},
		Printer:   self.Printer,
		PrintFlag: self.PrintFlag,
	}

	ret := newLeaf(self.Retry)

	for i := 0; i < self.Retry; i++ {
		answer := self.client.Query(q)
		ret.add(answer)
		if answer.Timeout() {
			self.Print("// retry")
			continue
		}
		break
	}

	return ret
}
