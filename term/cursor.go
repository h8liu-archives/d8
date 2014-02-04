package term

import (
	"errors"

	"github.com/h8liu/d8/client"
	"github.com/h8liu/d8/printer"
)

type Cursor interface {
	printer.Interface

	Error() error
	T(t Task) (*Branch, error)
	Q(q *client.Query) (*Leaf, error)
}

type cursor struct {
	*printer.Printer
	*Term // conveniently inherits the term options
	*stack
	nquery int
	e      error
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
	MaxDepth = 30
	MaxQuery = 500
)

var (
	errTooDeep        = errors.New("too deep")
	errTooManyQueries = errors.New("too many queries")
)

func (self *cursor) Error() error { return self.e }

func (self *cursor) Q(q *client.Query) (*Leaf, error) {
	if self.e != nil {
		return nil, self.e
	}

	if self.nquery >= MaxQuery {
		self.e = errTooManyQueries
		self.Printf("error %v", self.e)
		return nil, self.e
	}

	self.nquery++
	ret := self.q(q)
	self.TopAdd(ret)
	return ret, self.e
}

func (self *cursor) T(t Task) (*Branch, error) {
	if self.e != nil {
		return nil, self.e
	}

	if self.Len() >= MaxDepth {
		self.e = errTooDeep
		self.Printf("error %v", self.e)
		return nil, self.e
	}

	ret := newBranch(t)
	self.TopAdd(ret)
	self.Push(ret)

	t.Run(self)

	b := self.Pop()
	if b != ret || b.Task != t {
		panic("bug")
	}

	return ret, self.e
}

func (self *cursor) q(q *client.Query) *Leaf {
	qp := &client.QueryPrinter{
		Query:     q,
		Printer:   self.Printer,
		PrintFlag: self.PrintFlag,
	}

	ret := newLeaf(self.Retry)

	for i := 0; i < self.Retry; i++ {
		if i > 0 {
			self.Print("// retry")
		}
		answer := self.client.Query(qp)
		ret.add(answer)
		if answer.Timeout() {
			continue
		}
		break
	}

	return ret
}
