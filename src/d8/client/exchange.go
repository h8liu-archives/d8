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
	self.printSend(p)
	self.printRecv(p)
}

func (self *Exchange) printSend(p *printer.Printer) {
	p.Printf("%s {", self.Query.String())
	p.ShiftIn()

	switch self.Query.PrintFlag {
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

func (self *Exchange) printRecv(p *printer.Printer) {
	switch self.Query.PrintFlag {
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
