package printer

import (
	"bytes"
	"fmt"
	"io"
)

type Printer struct {
	Prefix string
	Indent string
	Shift  int
	Writer io.Writer
	Error  error
}

var _ Interface = new(Printer)

// Returns a new printer that writes to w
// if w is nil, than all prints to the printer will be noops
func New(w io.Writer) *Printer {
	if w == nil {
		w = noop
	}

	return &Printer{
		Indent: "    ",
		Writer: w,
	}
}

func (self *Printer) p(n *int, a ...interface{}) {
	if self.Error != nil {
		return
	}

	i, e := fmt.Fprint(self.Writer, a...)
	self.Error = e
	*n += i
}

func (self *Printer) pln(n *int, a ...interface{}) {
	if self.Error != nil {
		return
	}

	i, e := fmt.Fprintln(self.Writer, a...)
	self.Error = e
	*n += i
}

func (self *Printer) pf(n *int, format string, a ...interface{}) {
	if self.Error != nil {
		return
	}

	i, e := fmt.Fprintf(self.Writer, format, a...)
	self.Error = e
	*n += i
}

func (self *Printer) pre(n *int) {
	self.p(n, self.Prefix)
	for i := 0; i < self.Shift; i++ {
		self.p(n, self.Indent)
	}
}

func (self *Printer) Print(a ...interface{}) (int, error) {
	n := 0
	self.pre(&n)
	self.p(&n, a...)
	self.pln(&n)

	return n, self.Error
}

func (self *Printer) Println(a ...interface{}) (int, error) {
	n := 0
	self.pre(&n)
	self.pln(&n, a...)

	return n, self.Error
}

func (self *Printer) Printf(format string, a ...interface{}) (int, error) {
	n := 0
	self.pre(&n)
	self.pf(&n, format, a...)
	self.pln(&n)

	return n, self.Error
}

func (self *Printer) ShiftIn() {
	self.Shift++
}

func (self *Printer) ShiftOut(a ...interface{}) {
	if self.Shift == 0 {
		panic("shift already left most")
	}
	self.Shift--

	if len(a) > 0 {
		self.Print(a...)
	}
}

func String(p Printable) string {
	buf := new(bytes.Buffer)
	dev := New(buf)
	p.PrintTo(dev)
	return buf.String()
}
