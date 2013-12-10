package term

import (
	"errors"

	"d8/client"
	"printer"
)

type Cursor interface {
	printer.Interface

	T(t Task) (*Branch, error)
	Q(q *client.Query) (*Leaf, error)
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
	MaxQuery = 300
)

var (
	errTooDeep        = errors.New("too deep")
	errTooManyQueries = errors.New("too many queries")
)

func (self *cursor) Q(q *client.Query) (*Leaf, error) {
	if self.nquery >= MaxQuery {
		return nil, errTooManyQueries
	}

	self.nquery++
	ret := self.q(q)
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

func (self *cursor) q(q *client.Query) *Leaf {
	qp := &client.QueryPrinter{
		Query:     q,
		Printer:   self.Printer,
		PrintFlag: self.PrintFlag,
	}

	ret := newLeaf(self.Retry)

	for i := 0; i < self.Retry; i++ {
		answer := self.client.Query(qp)
		ret.add(answer)
		if answer.Timeout() {
			self.Print("// retry")
			continue
		}
		break
	}

	return ret
}
