package packet

import (
	"bytes"
	"printer"
)

type Section []*RR

func (self Section) LenU16() uint16 {
	if self == nil {
		return 0
	}

	if len(self) > 0xffff {
		panic("too many rrs")
	}

	return uint16(len(self))
}

func (self Section) unpack(in *bytes.Reader, p []byte) error {
	var e error
	for i, _ := range self {
		self[i], e = unpackRR(in, p)
		if e != nil {
			return e
		}
	}

	return nil
}

func (self Section) PrintTo(p *printer.Printer) {
	for _, rr := range self {
		rr.PrintTo(p)
	}
}

func (self Section) PrintNameTo(p *printer.Printer, name string) {
	p.Printf("%s {\n", name)
	p.ShiftIn()
	self.PrintTo(p)
	p.ShiftOut()
	p.Println("}")
}
