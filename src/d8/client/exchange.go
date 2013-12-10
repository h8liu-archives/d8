package client

import (
	"fmt"
	"printer"
)

type Exchange struct {
	Query     *Query
	Send      *Message
	Recv      *Message
	Error     error
	PrintFlag int
}

func (self *Exchange) PrintTo(p *printer.Printer) {
	self.printSend(p)
	self.printRecv(p)
}

func (self *Exchange) printSend(p *printer.Printer) {
	p.Printf("%s {", self.Query.String())
	p.ShiftIn()

	switch self.PrintFlag {
	case PrintAll:
		p.Print("send {")
		p.ShiftIn()
		self.Send.PrintTo(p)
		p.ShiftOut()
		p.Print("}")
	case PrintReply:
		// do nothing
	default:
		panic("unknown print flag")
	}
}

func (self *Exchange) printTimeTaken(p *printer.Printer) {
	d := self.Recv.Timestamp.Sub(self.Send.Timestamp)
	n := d.Nanoseconds()
	var s string
	if n < 1e3 {
		s = fmt.Sprintf("%dns", n)
	} else if n < 1e6 {
		s = fmt.Sprintf("%.1fus", float64(n)/1e3)
	} else if n < 1e9 {
		s = fmt.Sprintf("%.2fms", float64(n)/1e6)
	} else {
		s = fmt.Sprintf("%.3fs", float64(n)/1e9)
	}

	p.Printf("(in %v)", s)
}

func (self *Exchange) printRecv(p *printer.Printer) {
	switch self.PrintFlag {
	case PrintAll:
		if self.Recv != nil {
			p.Print("recv {")
			p.ShiftIn()
			self.Recv.PrintTo(p)
			p.ShiftOut()
			p.Print("}")
		}

		if self.Error != nil {
			p.Printf("error %v", self.Error)
		}
	case PrintReply:
		if self.Recv != nil {
			self.Recv.Packet.PrintTo(p)
			self.printTimeTaken(p)
		}
		if self.Error != nil {
			p.Printf("error %v", self.Error)
		}
	default:
		panic("unknown print flag")
	}

	p.ShiftOut()
	p.Print("}")
}

func (self *Exchange) String() string {
	return printer.String(self)
}

func (self *Exchange) Timeout() bool {
	return self.Error == errTimeout
}
