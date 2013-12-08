package client

import (
	"printer"
)

type Exchange struct {
	Query *Query
	Send  *Message
	Recv  *Message
	Error error
}

func (self *Exchange) PrintTo(p *printer.Printer) {
	p.Printf("%s {", self.Query.String())
	p.ShiftIn()

	p.Print("send {")
	p.ShiftIn()
	self.Send.PrintTo(p)
	p.ShiftOut()
	p.Print("}")

	if self.Recv != nil {
		p.Print("recv {")
		p.ShiftIn()
		self.Send.PrintTo(p)
		p.ShiftOut()
		p.Print("}")
	}

	if self.Error != nil {
		p.Printf("error %v", self.Error)
	}

	p.ShiftOut()
	p.Print("}")
}

func (self *Exchange) String() string {
	return printer.String(self)
}
